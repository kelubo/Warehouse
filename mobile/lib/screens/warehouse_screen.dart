import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

import '../models/models.dart';
import '../services/warehouse_service.dart';
import 'shelf_screen.dart';

class WarehouseScreen extends StatefulWidget {
  final Warehouse warehouse;

  const WarehouseScreen({super.key, required this.warehouse});

  @override
  State<WarehouseScreen> createState() => _WarehouseScreenState();
}

class _WarehouseScreenState extends State<WarehouseScreen> {
  @override
  void initState() {
    super.initState();
    _loadShelves();
  }

  Future<void> _loadShelves() async {
    await Provider.of<WarehouseService>(context, listen: false).loadShelves(widget.warehouse.id);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.warehouse.name),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () => _showEditDialog(context),
          ),
          IconButton(
            icon: const Icon(Icons.delete),
            onPressed: () => _deleteWarehouse(context),
          ),
        ],
      ),
      body: Consumer<WarehouseService>(
        builder: (context, service, child) {
          if (service.shelves.isEmpty) {
            return const Center(child: Text('暂无货架'));
          }
          return ListView.builder(
            itemCount: service.shelves.length,
            itemBuilder: (context, index) {
              final shelf = service.shelves[index];
              return Slidable(
                actionPane: const SlidableDrawerActionPane(),
                secondaryActions: [
                  IconSlideAction(
                    caption: '删除',
                    color: Colors.red,
                    icon: Icons.delete,
                    onTap: () => _deleteShelf(context, shelf.id),
                  ),
                ],
                child: ListTile(
                  title: Text(shelf.name),
                  subtitle: Text('${shelf.rows}层 x ${shelf.columns}列'),
                  trailing: const Icon(Icons.arrow_forward),
                  onTap: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => ShelfScreen(shelf: shelf),
                      ),
                    );
                  },
                ),
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddShelfDialog(context),
        child: const Icon(Icons.add),
      ),
    );
  }

  void _showAddShelfDialog(BuildContext context) {
    final nameController = TextEditingController();
    final rowsController = TextEditingController(text: '5');
    final colsController = TextEditingController(text: '5');

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('新建货架'),
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
              await Provider.of<WarehouseService>(context, listen: false).addShelf(
                widget.warehouse.id,
                nameController.text,
                int.parse(rowsController.text),
                int.parse(colsController.text),
              );
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showEditDialog(BuildContext context) {
    final nameController = TextEditingController(text: widget.warehouse.name);
    final descController = TextEditingController(text: widget.warehouse.description);
    final locationController = TextEditingController(text: widget.warehouse.location);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('编辑仓库'),
        content: SingleChildScrollView(
          child: Column(
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(labelText: '仓库名称'),
              ),
              TextField(
                controller: descController,
                decoration: const InputDecoration(labelText: '描述'),
              ),
              TextField(
                controller: locationController,
                decoration: const InputDecoration(labelText: '位置'),
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
              await Provider.of<WarehouseService>(context, listen: false).updateWarehouse(
                widget.warehouse.copyWith(
                  name: nameController.text,
                  description: descController.text,
                  location: locationController.text,
                ),
              );
              Navigator.pop(context);
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _deleteWarehouse(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认删除'),
        content: const Text('确定要删除这个仓库吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              await Provider.of<WarehouseService>(context, listen: false).removeWarehouse(widget.warehouse.id);
              Navigator.pop(context);
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _deleteShelf(BuildContext context, String shelfId) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认删除'),
        content: const Text('确定要删除这个货架吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              await Provider.of<WarehouseService>(context, listen: false).removeShelf(shelfId);
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }
}
