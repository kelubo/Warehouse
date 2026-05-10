import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/models.dart';
import '../services/warehouse_service.dart';

class BoxDetailScreen extends StatelessWidget {
  final Box box;

  const BoxDetailScreen({super.key, required this.box});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('箱子详情'),
        actions: [
          IconButton(
            icon: const Icon(Icons.delete),
            onPressed: () => _deleteBox(context),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildInfoRow('位置', '货架 ${box.shelfId.substring(0, 8)} - 行${box.row}列${box.column}'),
            _buildInfoRow('编号', box.code),
            _buildInfoRow('数量', box.quantity.toString()),
            _buildInfoRow('内容', box.contents.isNotEmpty ? box.contents : '无'),
            _buildInfoRow('创建时间', box.createdAt.toString()),
            _buildInfoRow('更新时间', box.updatedAt.toString()),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showEditDialog(context),
        child: const Icon(Icons.edit),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label, style: const TextStyle(color: Colors.grey)),
          Text(value, style: const TextStyle(fontSize: 18)),
        ],
      ),
    );
  }

  void _showEditDialog(BuildContext context) {
    final contentsController = TextEditingController(text: box.contents);
    final quantityController = TextEditingController(text: box.quantity.toString());

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('编辑箱子'),
        content: SingleChildScrollView(
          child: Column(
            children: [
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

  void _deleteBox(BuildContext context) {
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
              await Provider.of<WarehouseService>(context, listen: false).removeBox(box.id);
              Navigator.pop(context);
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }
}
