<?php

/**
 * Box packing (3D bin packing, knapsack problem).
 *
 * @author Doug Wright
 */
declare(strict_types=1);

namespace DVDoug\BoxPacker;

use DVDoug\BoxPacker\Test\TestBox;
use DVDoug\BoxPacker\Test\TestItem;
use PHPUnit\Framework\Attributes\CoversClass;
use PHPUnit\Framework\TestCase;

use function json_encode;

#[CoversClass(PackedBoxList::class)]
class PackedBoxListTest extends TestCase
{
    /**
     * Test that inserting individually correctly works.
     */
    public function testInsertAndCount(): void
    {
        $box = new TestBox('Box', 10, 10, 10, 0, 10, 10, 10, 100);
        $itemA = new TestItem('Item A', 5, 10, 10, 10, Rotation::BestFit);
        $itemB = new TestItem('Item B', 5, 10, 10, 20, Rotation::BestFit);

        $packedItemA = new PackedItem($itemA, 0, 0, 0, 5, 10, 10);
        $packedItemB = new PackedItem($itemB, 0, 0, 0, 5, 10, 10);

        $packedItemListA = new PackedItemList();
        $packedItemListA->insert($packedItemA);
        $packedBoxA = new PackedBox($box, $packedItemListA);

        $packedItemListB = new PackedItemList();
        $packedItemListB->insert($packedItemB);
        $packedBoxB = new PackedBox($box, $packedItemListB);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBoxA);
        $packedBoxList->insert($packedBoxB);

        self::assertCount(2, $packedBoxList);
    }

    /**
     * Test that inserting in bulk correctly works.
     */
    public function testInsertFromArrayAndCount(): void
    {
        $box = new TestBox('Box', 10, 10, 10, 0, 10, 10, 10, 100);
        $itemA = new TestItem('Item A', 5, 10, 10, 10, Rotation::BestFit);
        $itemB = new TestItem('Item B', 5, 10, 10, 20, Rotation::BestFit);

        $packedItemA = new PackedItem($itemA, 0, 0, 0, 5, 10, 10);
        $packedItemB = new PackedItem($itemB, 0, 0, 0, 5, 10, 10);

        $packedItemListA = new PackedItemList();
        $packedItemListA->insert($packedItemA);
        $packedBoxA = new PackedBox($box, $packedItemListA);

        $packedItemListB = new PackedItemList();
        $packedItemListB->insert($packedItemB);
        $packedBoxB = new PackedBox($box, $packedItemListB);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insertFromArray([$packedBoxA, $packedBoxB]);

        self::assertCount(2, $packedBoxList);
    }

    /**
     * Test we can peek at the "top" (next) item in the list.
     */
    public function testTop(): void
    {
        $box = new TestBox('Box', 10, 10, 10, 0, 10, 10, 10, 100);
        $itemA = new TestItem('Item A', 5, 10, 10, 10, Rotation::BestFit);
        $itemB = new TestItem('Item B', 5, 10, 10, 20, Rotation::BestFit);

        $packedItemA = new PackedItem($itemA, 0, 0, 0, 5, 10, 10);
        $packedItemB = new PackedItem($itemB, 0, 0, 0, 5, 10, 10);

        $packedItemListA = new PackedItemList();
        $packedItemListA->insert($packedItemA);
        $packedBoxA = new PackedBox($box, $packedItemListA);

        $packedItemListB = new PackedItemList();
        $packedItemListB->insert($packedItemB);
        $packedBoxB = new PackedBox($box, $packedItemListB);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBoxA);
        $packedBoxList->insert($packedBoxB);

        self::assertEquals($packedBoxA, $packedBoxList->top());
    }

    /**
     * Test that volume utilisation is correctly calculated.
     */
    public function testVolumeUtilisation(): void
    {
        $box = new TestBox('Box', 10, 10, 10, 0, 10, 10, 10, 10);
        $item = new TestItem('Item', 5, 10, 10, 10, Rotation::BestFit);

        $packedItem = new PackedItem($item, 0, 0, 0, 5, 10, 10);

        $packedItemList = new PackedItemList();
        $packedItemList->insert($packedItem);

        $packedBox = new PackedBox($box, $packedItemList);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBox);

        self::assertEquals(50, $packedBoxList->getVolumeUtilisation());
    }

    /**
     * Test that weight variance is correctly calculated.
     */
    public function testWeightVariance(): void
    {
        $box = new TestBox('Box', 10, 10, 10, 0, 10, 10, 10, 100);
        $itemA = new TestItem('Item A', 5, 10, 10, 10, Rotation::BestFit);
        $itemB = new TestItem('Item B', 5, 10, 10, 20, Rotation::BestFit);

        $packedItemA = new PackedItem($itemA, 0, 0, 0, 5, 10, 10);
        $packedItemB = new PackedItem($itemB, 0, 0, 0, 5, 10, 10);

        $packedItemListA = new PackedItemList();
        $packedItemListA->insert($packedItemA);
        $packedBoxA = new PackedBox($box, $packedItemListA);

        $packedItemListB = new PackedItemList();
        $packedItemListB->insert($packedItemB);
        $packedBoxB = new PackedBox($box, $packedItemListB);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBoxA);
        $packedBoxList->insert($packedBoxB);

        self::assertEquals(25, $packedBoxList->getWeightVariance());
    }

    /**
     * Test that mean weight is correctly calculated.
     */
    public function testMeanWeight(): void
    {
        $box = new TestBox('Box', 10, 10, 10, 0, 10, 10, 10, 100);
        $itemA = new TestItem('Item A', 5, 10, 10, 10, Rotation::BestFit);
        $itemB = new TestItem('Item B', 5, 10, 10, 20, Rotation::BestFit);

        $packedItemA = new PackedItem($itemA, 0, 0, 0, 5, 10, 10);
        $packedItemB = new PackedItem($itemB, 0, 0, 0, 5, 10, 10);

        $packedItemListA = new PackedItemList();
        $packedItemListA->insert($packedItemA);
        $packedBoxA = new PackedBox($box, $packedItemListA);

        $packedItemListB = new PackedItemList();
        $packedItemListB->insert($packedItemB);
        $packedBoxB = new PackedBox($box, $packedItemListB);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBoxA);
        $packedBoxList->insert($packedBoxB);

        self::assertEquals(15, $packedBoxList->getMeanWeight());
    }

    /**
     * Test JSON representation.
     */
    public function testJsonSerialize(): void
    {
        $box = new TestBox('Box', 10, 10, 20, 10, 10, 10, 20, 10);
        $item = new OrientatedItem(new TestItem('Item', 4, 10, 10, 10, Rotation::KeepFlat), 4, 10, 10);

        $boxItems = new PackedItemList();
        $boxItems->insert(PackedItem::fromOrientatedItem($item, 0, 0, 0));

        $packedBox = new PackedBox($box, $boxItems);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBox);

        self::assertJsonStringEqualsJsonString('[{"box":{"reference":"Box","innerWidth":10,"innerLength":10,"innerDepth":20,"emptyWeight": 10,"maxWeight": 10},"items":[{"x":0,"y":0,"z":0,"width":4,"length":10,"depth":10,"item":{"description":"Item","width":4,"length":10,"depth":10,"allowedRotation":2,"weight": 10}}]}]', json_encode($packedBoxList));
    }

    /**
     * Test visualisation URL.
     */
    public function testVisualisationURL(): void
    {
        $box = new TestBox('Box', 10, 10, 20, 10, 10, 10, 20, 10);
        $item = new OrientatedItem(new TestItem('Item', 4, 10, 10, 10, Rotation::KeepFlat), 4, 10, 10);

        $boxItems = new PackedItemList();
        $boxItems->insert(PackedItem::fromOrientatedItem($item, 0, 0, 0));

        $packedBox = new PackedBox($box, $boxItems);

        $packedBoxList = new PackedBoxList();
        $packedBoxList->insert($packedBox);

        self::assertEquals('https://boxpacker.io/en/master/visualiser.html?packing=%7B%22items%22%3A%5B%5B%22Item%22%2C4%2C10%2C10%5D%5D%2C%22boxes%22%3A%5B%5B%22Box%22%2C10%2C10%2C20%2C%5B%5B0%2C0%2C0%2C0%2C4%2C10%2C10%5D%5D%5D%5D%7D', $packedBoxList->generateVisualisationURL());
    }
}
