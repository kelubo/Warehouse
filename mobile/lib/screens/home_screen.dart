import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../services/warehouse_service.dart';
import '../services/sync_service.dart';
import 'warehouse_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    await Provider.of<WarehouseService>(context, listen: false).loadData();
    await Provider.of<SyncService>(context, listen: false).checkPendingChanges();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('仓库管理系统'),
        actions: [
          Consumer<SyncService>(
            builder: (context, syncService, child) {
              return IconButton(
                icon: syncService.isSyncing
                    ? const CircularProgressIndicator()
                    : const Icon(Icons.sync),
                onPressed: () async {
                  try {
                    await syncService.sync();
                    await _loadData();
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('同步成功')),
                    );
                  } catch (e) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('同步失败: $e')),
                    );
                  }
                },
              );
            },
          ),
        ],
      ),
      body: Consumer<WarehouseService>(
        builder: (context, service, child) {
          if (service.warehouses.isEmpty) {
            return const Center(child: Text('暂无仓库'));
          }
          return ListView.builder(
            itemCount: service.warehouses.length,
            itemBuilder: (context, index) {
              final warehouse = service.warehouses[index];
              return ListTile(
                title: Text(warehouse.name),
                subtitle: Text(warehouse.location),
                trailing: const Icon(Icons.arrow_forward),
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => WarehouseScreen(warehouse: warehouse),
                    ),
                  );
                },
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddWarehouseDialog(context),
        child: const Icon(Icons.add),
      ),
    );
  }

  void _showAddWarehouseDialog(BuildContext context) {
    final nameController = TextEditingController();
    final descController = TextEditingController();
    final locationController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('新建仓库'),
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
              await Provider.of<WarehouseService>(context, listen: false).addWarehouse(
                nameController.text,
                descController.text,
                locationController.text,
              );
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }
}
