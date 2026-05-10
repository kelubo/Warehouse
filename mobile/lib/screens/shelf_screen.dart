import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

import '../models/models.dart';
import '../services/warehouse_service.dart';

class ShelfScreen extends StatefulWidget {
  final Shelf shelf;

  const ShelfScreen({super.key, required this.shelf});

  @override
  State<ShelfScreen> createState() => _ShelfScreenState();
}

class _ShelfScreenState extends State<ShelfScreen> {
  @override
  void initState() {
    super.initState();
    _loadBoxes();
  }

  Future<void> _loadBoxes() async {
    await Provider.of<WarehouseService>(context, listen: false).loadBoxes(widget.shelf.id);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.shelf.name),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () => _showEditDialog(context),
          ),
        ],
      ),
      body: Consumer<WarehouseService>(
        builder: (context, service, child) {
          return Column(
            children: [
              _buildGrid(context, service.boxes),
              const SizedBox(height: 20),
              const Divider(),
              const Text('箱子列表'),
              Expanded(
                child: ListView.builder(
                  itemCount: service.boxes.length,
                  itemBuilder: (context, index) {
                    final box = service.boxes[index];
                    return Slidable(
                      actionPane: const SlidableDrawerActionPane(),
                      secondaryActions: [
                        IconSlideAction(
                          caption: '删除',
                          color: Colors.red,
                          icon: Icons.delete,
                          onTap: () => _deleteBox(context, box.id),
                        ),
                      ],
                      child: ListTile(
                        title: Text('${box.row}-${box.column}: ${box.code}'),
                        subtitle: Text('数量: ${box.quantity}'),
                        onTap: () => _showBoxDetail(context, box),
                      ),
                    );
                  },
                ),
              ),
            ],
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddBoxDialog(context),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildGrid(BuildContext context, List<Box> boxes) {
    return GridView.builder(
      shrinkWrap: true,
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: widget.shelf.columns,
        childAspectRatio: 1.0,
      ),
      itemCount: widget.shelf.rows * widget.shelf.columns,
      itemBuilder: (context, index) {
        final row = (index ~/ widget.shelf.columns) + 1;
        final col = (index % widget.shelf.columns) + 1;
        final box = boxes.firstWhere(
          (b) => b.row == row && b.column == col,
          orElse: () => null,
        );

        return Container(
          margin: const EdgeInsets.all(2),
          decoration: BoxDecoration(
            border: Border.all(color: Colors.grey),
            color: box != null ? Colors.blue[100] : Colors.grey[100],
          ),
          child: Center(
            child: box != null
                ? Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(box.code, style: const TextStyle(fontSize: 10)),
                      Text('${box.quantity}', style: const TextStyle(fontSize: 12)),
                    ],
                  )
                : Text('$row-$col', style: const TextStyle(color: Colors.grey)),
          ),
        );
      },
    );
  }

  void _showAddBoxDialog(BuildContext context) {
    final codeController = TextEditingController();
    final contentsController = TextEditingController();
    final quantityController = TextEditingController(text: '1');

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('新建箱子'),
        content: SingleChildScrollView(
          child: Column(
            children: [
              TextField(
                controller: codeController,
                decoration: const InputDecoration(labelText: '箱子编号'),
              ),
              TextField(
                controller: contentsController,
                decoration: const InputDecoration(labelText: '内容描述'),
              ),
              TextField(
                controller: quantityController,
                decoration: const InputDecoration(labelText: '数量'),
                keyboardType: TextInputType.number,
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              final service = Provider.of<WarehouseService>(context, listen: false);
              final nextPosition = _findNextPosition(service.boxes);
              await service.addBox(
                widget.shelf.id,
                nextPosition.$1,
                nextPosition.$2,
                codeController.text,
                contentsController.text,
                int.parse(quantityController.text),
              );
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  (int, int) _findNextPosition(List<Box> boxes) {
    for (int row = 1; row <= widget.shelf.rows; row++) {
      for (int col = 1; col <= widget.shelf.columns; col++) {
        if (!boxes.any((b) => b.row == row && b.column == col)) {
          return (row, col);
        }
      }
    }
    return (1, 1);
  }

  void _showBoxDetail(BuildContext context, Box box) {
    final contentsController = TextEditingController(text: box.contents);
    final quantityController = TextEditingController(text: box.quantity.toString());

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('位置: ${box.row}-${box.column}'),
        content: SingleChildScrollView(
          child: Column(
            children: [
              Text('编号: ${box.code}'),
              TextField(
                controller: contentsController,
                decoration: const InputDecoration(labelText: '内容描述'),
              ),
              TextField(
                controller: quantityController,
                decoration: const InputDecoration(labelText: '数量'),
                keyboardType: TextInputType.number,
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              await Provider.of<WarehouseService>(context, listen: false).updateBox(
                box.copyWith(
                  contents: contentsController.text,
                  quantity: int.parse(quantityController.text),
                ),
              );
              Navigator.pop(context);
            },
            child: const Text('保存'),
          ),
        ],
      ),
    );
  }

  void _showEditDialog(BuildContext context) {
    final nameController = TextEditingController(text: widget.shelf.name);
    final rowsController = TextEditingController(text: widget.shelf.rows.toString());
    final colsController = TextEditingController(text: widget.shelf.columns.toString());

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('编辑货架'),
        content: SingleChildScrollView(
          child: Column(
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(labelText: '货架名称'),
              ),
              TextField(
                controller: rowsController,
                decoration: const InputDecoration(labelText: '层数'),
                keyboardType: TextInputType.number,
              ),
              TextField(
                controller: colsController,
                decoration: const InputDecoration(labelText: '列数'),
                keyboardType: TextInputType.number,
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              await Provider.of<WarehouseService>(context, listen: false).updateShelf(
                widget.shelf.copyWith(
                  name: nameController.text,
                  rows: int.parse(rowsController.text),
                  columns: int.parse(colsController.text),
                ),
              );
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _deleteBox(BuildContext context, String boxId) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认删除'),
        content: const Text('确定要删除这个箱子吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              await Provider.of<WarehouseService>(context, listen: false).removeBox(boxId);
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }
}
