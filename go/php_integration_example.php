<?php

/**
 * Example PHP integration with Go lookahead library
 *
 * This file demonstrates how to integrate the Go-based lookahead
 * optimization into the existing BoxPacker PHP codebase.
 */

namespace DVDoug\BoxPacker;

use FFI;

/**
 * Go-accelerated version of OrientatedItemSorter
 *
 * Drop-in replacement that uses the Go library for lookahead calculations
 */
class OrientatedItemSorterGo
{
    private FFI $ffi;
    private bool $useGo = true;

    public function __construct(
        private readonly OrientatedItemFactory $orientatedItemFactory,
        private readonly bool $singlePassMode,
        private readonly int $widthLeft,
        private readonly int $lengthLeft,
        private readonly int $depthLeft,
        private readonly ItemList $nextItems,
        private readonly int $rowLength,
        private readonly int $x,
        private readonly int $y,
        private readonly int $z,
        private readonly PackedItemList $prevPackedItemList,
        private readonly \Psr\Log\LoggerInterface $logger
    ) {
        $this->initFFI();
    }

    private function initFFI(): void
    {
        try {
            $libPath = __DIR__ . '/libboxpacker.so';

            // On macOS use .dylib, on Windows use .dll
            if (PHP_OS_FAMILY === 'Darwin') {
                $libPath = __DIR__ . '/libboxpacker.dylib';
            } elseif (PHP_OS_FAMILY === 'Windows') {
                $libPath = __DIR__ . '/libboxpacker.dll';
            }

            if (!file_exists($libPath)) {
                $this->logger->warning("Go library not found at {$libPath}, falling back to PHP implementation");
                $this->useGo = false;
                return;
            }

            $this->ffi = FFI::cdef("
                typedef struct {
                    int width;
                    int length;
                    int depth;
                    int weight;
                    int rotation;
                } CItem;

                typedef struct {
                    int innerWidth;
                    int innerLength;
                    int innerDepth;
                    int maxWeight;
                } CBox;

                typedef struct {
                    int width;
                    int length;
                    int depth;
                    int surfaceFootprint;
                } COrientatedItem;

                int CalculateLookaheadFFI(
                    int prevItemWidth, int prevItemLength, int prevItemDepth,
                    CItem* items, int itemCount,
                    int widthLeft, int lengthLeft, int depthLeft, int rowLength,
                    int maxLookahead
                );

                int GetBestOrientationFFI(
                    CItem* item,
                    CItem* nextItems,
                    int nextItemCount,
                    int widthLeft, int lengthLeft, int depthLeft,
                    int rowLength,
                    int packedWeight,
                    CBox* box,
                    COrientatedItem* resultOrientation
                );

                void ClearCacheFFI();
                int GetCacheSizeFFI();
            ", $libPath);

            $this->logger->info("Go lookahead library loaded successfully");
        } catch (\Throwable $e) {
            $this->logger->warning("Failed to load Go library: {$e->getMessage()}, falling back to PHP");
            $this->useGo = false;
        }
    }

    public function __invoke(OrientatedItem $a, OrientatedItem $b): int
    {
        // If Go library is not available, fall back to original PHP implementation
        if (!$this->useGo) {
            return $this->phpCompare($a, $b);
        }

        try {
            return $this->goCompare($a, $b);
        } catch (\Throwable $e) {
            $this->logger->warning("Go comparison failed: {$e->getMessage()}, falling back to PHP");
            $this->useGo = false;
            return $this->phpCompare($a, $b);
        }
    }

    private function goCompare(OrientatedItem $a, OrientatedItem $b): int
    {
        // Exact fit checks (keep in PHP as they're fast)
        $orientationAWidthLeft = $this->widthLeft - $a->width;
        $orientationBWidthLeft = $this->widthLeft - $b->width;
        $widthDecider = $this->exactFitDecider($orientationAWidthLeft, $orientationBWidthLeft);
        if ($widthDecider !== 0) {
            return $widthDecider;
        }

        $orientationALengthLeft = $this->lengthLeft - $a->length;
        $orientationBLengthLeft = $this->lengthLeft - $b->length;
        $lengthDecider = $this->exactFitDecider($orientationALengthLeft, $orientationBLengthLeft);
        if ($lengthDecider !== 0) {
            return $lengthDecider;
        }

        $orientationADepthLeft = $this->depthLeft - $a->depth;
        $orientationBDepthLeft = $this->depthLeft - $b->depth;
        $depthDecider = $this->exactFitDecider($orientationADepthLeft, $orientationBDepthLeft);
        if ($depthDecider !== 0) {
            return $depthDecider;
        }

        // Use Go for lookahead calculation
        if ($this->nextItems->count() > 0) {
            $followingItemDecider = $this->goLookAheadDecider($a, $b, $orientationAWidthLeft, $orientationBWidthLeft);
            if ($followingItemDecider !== 0) {
                return $followingItemDecider;
            }
        }

        // Fallback to simple comparison
        $orientationAMinGap = min($orientationAWidthLeft, $orientationALengthLeft);
        $orientationBMinGap = min($orientationBWidthLeft, $orientationBLengthLeft);

        return $orientationAMinGap <=> $orientationBMinGap ?: $a->surfaceFootprint <=> $b->surfaceFootprint;
    }

    private function goLookAheadDecider(OrientatedItem $a, OrientatedItem $b, int $orientationAWidthLeft, int $orientationBWidthLeft): int
    {
        // Convert next items to C struct array
        $nextItemsArray = $this->convertItemsToC(
            iterator_to_array($this->nextItems->topN(8))
        );

        // Calculate lookahead for orientation A
        $additionalPackedA = $this->ffi->CalculateLookaheadFFI(
            $a->width,
            $a->length,
            $a->depth,
            $nextItemsArray,
            min(8, $this->nextItems->count()),
            $this->widthLeft,
            $this->lengthLeft,
            $this->depthLeft,
            $this->rowLength,
            8
        );

        // Calculate lookahead for orientation B
        $additionalPackedB = $this->ffi->CalculateLookaheadFFI(
            $b->width,
            $b->length,
            $b->depth,
            $nextItemsArray,
            min(8, $this->nextItems->count()),
            $this->widthLeft,
            $this->lengthLeft,
            $this->depthLeft,
            $this->rowLength,
            8
        );

        return $additionalPackedB <=> $additionalPackedA ?: 0;
    }

    private function convertItemsToC(array $items): FFI\CData
    {
        $count = count($items);
        $cItems = $this->ffi->new("CItem[$count]");

        foreach ($items as $i => $item) {
            $cItems[$i]->width = $item->getWidth();
            $cItems[$i]->length = $item->getLength();
            $cItems[$i]->depth = $item->getDepth();
            $cItems[$i]->weight = $item->getWeight();
            $cItems[$i]->rotation = $item->getAllowedRotation()->value;
        }

        return $cItems;
    }

    /**
     * Fallback to original PHP implementation
     */
    private function phpCompare(OrientatedItem $a, OrientatedItem $b): int
    {
        // Original OrientatedItemSorter logic
        $sorter = new \DVDoug\BoxPacker\OrientatedItemSorter(
            $this->orientatedItemFactory,
            $this->singlePassMode,
            $this->widthLeft,
            $this->lengthLeft,
            $this->depthLeft,
            $this->nextItems,
            $this->rowLength,
            $this->x,
            $this->y,
            $this->z,
            $this->prevPackedItemList,
            $this->logger
        );

        return $sorter($a, $b);
    }

    private function exactFitDecider(int $dimensionALeft, int $dimensionBLeft): int
    {
        if ($dimensionALeft === 0 && $dimensionBLeft > 0) {
            return -1;
        }

        if ($dimensionALeft > 0 && $dimensionBLeft === 0) {
            return 1;
        }

        return 0;
    }

    /**
     * Clear the Go cache (call between packing jobs)
     */
    public function clearCache(): void
    {
        if ($this->useGo) {
            $this->ffi->ClearCacheFFI();
        }
    }

    /**
     * Get cache statistics
     */
    public function getCacheSize(): int
    {
        if ($this->useGo) {
            return $this->ffi->GetCacheSizeFFI();
        }
        return 0;
    }
}

/**
 * Usage example:
 */
function exampleUsage()
{
    // In OrientatedItemFactory.php, replace:
    // $sorter = new OrientatedItemSorter(...);

    // With:
    // $sorter = new OrientatedItemSorterGo(...);

    // The rest of the code remains unchanged!
}
