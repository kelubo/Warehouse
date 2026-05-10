const API_BASE = '/api/v1';
let currentPage = 'dashboard';
let token = localStorage.getItem('token');

function showLoginPage() {
    document.getElementById('app').innerHTML = `
        <div class="login-container">
            <div class="login-box">
                <h2>🏠 家庭储物管理</h2>
                <form id="loginForm">
                    <div class="form-group">
                        <label>用户名</label>
                        <input type="text" id="loginUsername" required placeholder="请输入用户名">
                    </div>
                    <div class="form-group">
                        <label>密码</label>
                        <input type="password" id="loginPassword" required placeholder="请输入密码">
                    </div>
                    <button type="submit" class="btn btn-primary">登录</button>
                </form>
                <div class="small-link">
                    <a href="#" onclick="showRegisterForm()">没有账户？立即注册</a>
                </div>
            </div>
        </div>
    `;
    document.getElementById('loginForm').addEventListener('submit', handleLogin);
}

function showRegisterForm() {
    document.getElementById('app').innerHTML = `
        <div class="login-container">
            <div class="login-box">
                <h2>🏠 注册账户</h2>
                <form id="registerForm">
                    <div class="form-group">
                        <label>用户名</label>
                        <input type="text" id="registerUsername" required placeholder="请输入用户名">
                    </div>
                    <div class="form-group">
                        <label>邮箱</label>
                        <input type="email" id="registerEmail" required placeholder="请输入邮箱">
                    </div>
                    <div class="form-group">
                        <label>密码</label>
                        <input type="password" id="registerPassword" required placeholder="请输入密码">
                    </div>
                    <button type="submit" class="btn btn-primary">注册</button>
                </form>
                <div class="small-link">
                    <a href="#" onclick="showLoginPage()">已有账户？立即登录</a>
                </div>
            </div>
        </div>
    `;
    document.getElementById('registerForm').addEventListener('submit', handleRegister);
}

function handleLogin(e) {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;
    fetch(`${API_BASE}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.token) {
            localStorage.setItem('token', data.token);
            token = data.token;
            showMainPage();
        } else {
            showNotification('登录失败：' + (data.message || '用户名或密码错误'), 'error');
        }
    })
    .catch(error => {
        showNotification('登录失败：网络错误', 'error');
    });
}

function handleRegister(e) {
    e.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;
    fetch(`${API_BASE}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.id) {
            showNotification('注册成功，请登录', 'success');
            setTimeout(showLoginPage, 2000);
        } else {
            showNotification('注册失败：' + (data.message || '未知错误'), 'error');
        }
    })
    .catch(error => {
        showNotification('注册失败：网络错误', 'error');
    });
}

function showMainPage() {
    document.getElementById('app').innerHTML = `
        <div class="main-container">
            <div class="sidebar">
                <div class="logo">🏠 家庭储物</div>
                <ul class="menu">
                    <li><a href="#" onclick="loadPage('dashboard')" id="menu-dashboard"><i>📊</i> 总览</a></li>
                    <li><a href="#" onclick="loadPage('products')" id="menu-products"><i>📦</i> 物品管理</a></li>
                    <li><a href="#" onclick="loadPage('warehouses')" id="menu-warehouses"><i>🏠</i> 储物空间</a></li>
                    <li><a href="#" onclick="loadPage('shelves')" id="menu-shelves"><i>📚</i> 置物架</a></li>
                    <li><a href="#" onclick="loadPage('boxes')" id="menu-boxes"><i>📦</i> 收纳盒</a></li>
                </ul>
                <div class="user-info">
                    <button class="logout-btn" onclick="logout()">退出登录</button>
                </div>
            </div>
            <div class="content">
                <div class="content-header">
                    <h1 id="page-title">总览</h1>
                </div>
                <div id="notification" class="notification"></div>
                <div id="page-content"></div>
            </div>
        </div>
    `;
    loadPage('dashboard');
}

function loadPage(page) {
    currentPage = page;
    document.querySelectorAll('.menu a').forEach(link => link.classList.remove('active'));
    document.getElementById('menu-' + page)?.classList.add('active');
    const titles = {
        dashboard: '总览',
        products: '物品管理',
        warehouses: '储物空间',
        shelves: '置物架',
        boxes: '收纳盒'
    };
    document.getElementById('page-title').textContent = titles[page];
    switch(page) {
        case 'dashboard': loadDashboard(); break;
        case 'products': loadProducts(); break;
        case 'warehouses': loadWarehouses(); break;
        case 'shelves': loadShelves(); break;
        case 'boxes': loadBoxes(); break;
    }
}

function loadDashboard() {
    Promise.all([
        fetch(`${API_BASE}/products`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/warehouses`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/shelves`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/boxes`, { headers: getAuthHeader() }).then(r => r.json())
    ])
    .then(([products, warehouses, shelves, boxes]) => {
        document.getElementById('page-content').innerHTML = `
            <div class="stats-grid">
                <div class="stat-card"><div class="number">${products.length}</div><div class="label">物品总数</div></div>
                <div class="stat-card"><div class="number">${warehouses.length}</div><div class="label">储物空间</div></div>
                <div class="stat-card"><div class="number">${shelves.length}</div><div class="label">置物架数</div></div>
                <div class="stat-card"><div class="number">${boxes.length}</div><div class="label">收纳盒数</div></div>
            </div>
            <div class="card">
                <div class="card-header">
                    <h2>最近添加的物品</h2>
                    <button class="btn btn-primary" onclick="loadPage('products')">查看全部</button>
                </div>
                <table>
                    <tr><th>名称</th><th>分类</th><th>存放位置</th></tr>
                    ${products.slice(-5).reverse().map(p => `<tr><td>${p.name}</td><td>${p.category || '-'}</td><td>${p.shelf_id ? '置物架' : p.box_id ? '收纳盒' : '未指定'}</td></tr>`).join('')}
                </table>
            </div>
        `;
    });
}

function loadProducts() {
    Promise.all([
        fetch(`${API_BASE}/products`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/shelves`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/boxes`, { headers: getAuthHeader() }).then(r => r.json())
    ])
    .then(([products, shelves, boxes]) => {
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>物品列表</h2>
                    <button class="btn btn-primary" onclick="showProductModal()">添加物品</button>
                </div>
                <table>
                    <tr><th>名称</th><th>分类</th><th>数量</th><th>存放位置</th><th>操作</th></tr>
                    ${products.map(p => {
                        const shelf = shelves.find(s => s.id === p.shelf_id);
                        const box = boxes.find(b => b.id === p.box_id);
                        let location = '未指定';
                        if (box) location = '收纳盒: ' + box.box_no;
                        else if (shelf && p.shelf_column > 0) location = shelf.name + ' (' + p.shelf_column + '-' + p.shelf_row + ')';
                        else if (shelf) location = shelf.name;
                        return '<tr><td>' + p.name + '</td><td>' + (p.category || '-') + '</td><td>' + (p.quantity || 1) + '</td><td>' + location + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="editProduct(' + JSON.stringify(p) + ')">编辑</button><button class="btn btn-danger" onclick="deleteProduct(\'' + p.id + '\')">删除</button></div></td></tr>';
                    }).join('')}
                </table>
            </div>
        `;
    });
}

function editProduct(product) {
    showProductModal(product);
}

function showProductModal(product = null) {
    const isEdit = product !== null;
    Promise.all([
        fetch(`${API_BASE}/shelves`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/boxes`, { headers: getAuthHeader() }).then(r => r.json())
    ])
    .then(([shelves, boxes]) => {
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()">
                    <h3>${isEdit ? '编辑物品' : '添加物品'}</h3>
                    <form id="productForm">
                        <div class="form-group"><label>物品名称</label><input type="text" id="productName" value="${product?.name || ''}" required></div>
                        <div class="form-row">
                            <div class="form-group"><label>分类</label><input type="text" id="productCategory" value="${product?.category || ''}" placeholder="如：衣物、书籍、电子产品"></div>
                            <div class="form-group"><label>数量</label><input type="number" id="productQuantity" value="${product?.quantity || 1}" min="1"></div>
                        </div>
                        <div class="form-group"><label>存放位置 - 置物架</label><select id="productShelf"><option value="">无</option>${shelves.map(s => '<option value="' + s.id + '"' + (product?.shelf_id === s.id ? ' selected' : '') + '>' + s.name + '</option>').join('')}</select></div>
                        <div class="form-group"><label>存放位置 - 收纳盒</label><select id="productBox"><option value="">无</option>${boxes.map(b => '<option value="' + b.id + '"' + (product?.box_id === b.id ? ' selected' : '') + '>' + b.box_no + '</option>').join('')}</select></div>
                        <div class="form-row" id="shelfPosition" ${product?.shelf_id ? '' : 'style="display:none"'}>
                            <div class="form-group"><label>列</label><input type="number" id="productShelfColumn" value="${product?.shelf_column || ''}" min="1"></div>
                            <div class="form-group"><label>层</label><input type="number" id="productShelfRow" value="${product?.shelf_row || ''}" min="1"></div>
                        </div>
                        <div class="form-group"><label>备注</label><textarea id="productDescription">${product?.description || ''}</textarea></div>
                        <input type="hidden" id="productId" value="${product?.id || ''}">
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                            <button type="submit" class="btn btn-primary">${isEdit ? '保存' : '添加'}</button>
                        </div>
                    </form>
                </div>
            </div>
        `;
        document.getElementById('productForm').addEventListener('submit', handleProductSubmit);
        document.getElementById('productShelf').addEventListener('change', function() {
            document.getElementById('shelfPosition').style.display = this.value ? '' : 'none';
        });
    });
}

function handleProductSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('productId').value;
    const data = {
        name: document.getElementById('productName').value,
        category: document.getElementById('productCategory').value,
        quantity: parseInt(document.getElementById('productQuantity').value) || 1,
        shelf_id: document.getElementById('productShelf').value,
        box_id: document.getElementById('productBox').value,
        shelf_column: parseInt(document.getElementById('productShelfColumn').value) || 0,
        shelf_row: parseInt(document.getElementById('productShelfRow').value) || 0,
        description: document.getElementById('productDescription').value
    };
    const url = id ? `${API_BASE}/products/${id}` : `${API_BASE}/products`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(() => { showNotification(id ? '物品更新成功' : '物品添加成功', 'success'); closeModal(); loadProducts(); })
    .catch(() => { showNotification('操作失败', 'error'); });
}

function deleteProduct(id) {
    if (!confirm('确定要删除这个物品吗？')) return;
    fetch(`${API_BASE}/products/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(() => { showNotification('物品删除成功', 'success'); loadProducts(); })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadWarehouses() {
    fetch(`${API_BASE}/warehouses`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(warehouses => {
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>储物空间列表</h2>
                    <button class="btn btn-primary" onclick="showWarehouseModal()">添加储物空间</button>
                </div>
                <table>
                    <tr><th>名称</th><th>位置</th><th>容量</th><th>描述</th><th>操作</th></tr>
                    ${warehouses.map(w => '<tr><td>' + w.name + '</td><td>' + (w.location || '-') + '</td><td>' + (w.capacity || '-') + '</td><td>' + (w.description || '-') + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="editWarehouse(' + JSON.stringify(w) + ')">编辑</button><button class="btn btn-danger" onclick="deleteWarehouse(\'' + w.id + '\')">删除</button></div></td></tr>').join('')}
                </table>
            </div>
        `;
    });
}

function editWarehouse(warehouse) {
    showWarehouseModal(warehouse);
}

function showWarehouseModal(warehouse = null) {
    const isEdit = warehouse !== null;
    document.getElementById('page-content').innerHTML += `
        <div class="modal-overlay active" onclick="closeModal()">
            <div class="modal" onclick="event.stopPropagation()">
                <h3>${isEdit ? '编辑储物空间' : '添加储物空间'}</h3>
                <form id="warehouseForm">
                    <div class="form-group"><label>名称</label><input type="text" id="warehouseName" value="${warehouse?.name || ''}" required></div>
                    <div class="form-group"><label>位置</label><input type="text" id="warehouseLocation" value="${warehouse?.location || ''}" placeholder="如：客厅、卧室、地下室"></div>
                    <div class="form-group"><label>容量</label><input type="number" id="warehouseCapacity" value="${warehouse && warehouse.capacity !== undefined ? warehouse.capacity : ''}" placeholder="可选"></div>
                    <div class="form-group"><label>描述</label><textarea id="warehouseDescription">${warehouse?.description || ''}</textarea></div>
                    <input type="hidden" id="warehouseId" value="${warehouse?.id || ''}">
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                        <button type="submit" class="btn btn-primary">${isEdit ? '保存' : '添加'}</button>
                    </div>
                </form>
            </div>
        </div>
    `;
    document.getElementById('warehouseForm').addEventListener('submit', handleWarehouseSubmit);
}

function handleWarehouseSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('warehouseId').value;
    const data = {
        name: document.getElementById('warehouseName').value,
        location: document.getElementById('warehouseLocation').value,
        capacity: parseInt(document.getElementById('warehouseCapacity').value) || 0,
        description: document.getElementById('warehouseDescription').value
    };
    const url = id ? `${API_BASE}/warehouses/${id}` : `${API_BASE}/warehouses`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(() => { showNotification(id ? '储物空间更新成功' : '储物空间添加成功', 'success'); closeModal(); loadWarehouses(); })
    .catch(() => { showNotification('操作失败', 'error'); });
}

function deleteWarehouse(id) {
    if (!confirm('确定要删除这个储物空间吗？')) return;
    fetch(`${API_BASE}/warehouses/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(() => { showNotification('储物空间删除成功', 'success'); loadWarehouses(); })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadShelves() {
    Promise.all([
        fetch(`${API_BASE}/shelves`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/warehouses`, { headers: getAuthHeader() }).then(r => r.json())
    ])
    .then(([shelves, warehouses]) => {
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>置物架列表</h2>
                    <button class="btn btn-primary" onclick="showShelfModal()">添加置物架</button>
                </div>
                <table>
                    <tr><th>名称</th><th>所属空间</th><th>列数</th><th>层数</th><th>操作</th></tr>
                    ${shelves.map(s => {
                        const warehouse = warehouses.find(w => w.id === s.warehouse_id);
                        return '<tr><td>' + s.name + '</td><td>' + (warehouse?.name || '-') + '</td><td>' + s.columns + '</td><td>' + s.rows + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="editShelf(' + JSON.stringify(s) + ')">编辑</button><button class="btn btn-danger" onclick="deleteShelf(\'' + s.id + '\')">删除</button></div></td></tr>';
                    }).join('')}
                </table>
            </div>
        `;
    });
}

function editShelf(shelf) {
    showShelfModal(shelf);
}

function showShelfModal(shelf = null) {
    const isEdit = shelf !== null;
    fetch(`${API_BASE}/warehouses`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(warehouses => {
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()">
                    <h3>${isEdit ? '编辑置物架' : '添加置物架'}</h3>
                    <form id="shelfForm">
                        <div class="form-group"><label>名称</label><input type="text" id="shelfName" value="${shelf?.name || ''}" required></div>
                        <div class="form-group"><label>所属储物空间</label><select id="shelfWarehouse"><option value="">请选择</option>${warehouses.map(w => '<option value="' + w.id + '"' + (shelf?.warehouse_id === w.id ? ' selected' : '') + '>' + w.name + '</option>').join('')}</select></div>
                        <div class="form-row">
                            <div class="form-group"><label>列数</label><input type="number" id="shelfColumns" value="${shelf?.columns || 1}" min="1" required></div>
                            <div class="form-group"><label>层数</label><input type="number" id="shelfRows" value="${shelf?.rows || 1}" min="1" required></div>
                        </div>
                        <input type="hidden" id="shelfId" value="${shelf?.id || ''}">
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                            <button type="submit" class="btn btn-primary">${isEdit ? '保存' : '添加'}</button>
                        </div>
                    </form>
                </div>
            </div>
        `;
        document.getElementById('shelfForm').addEventListener('submit', handleShelfSubmit);
    });
}

function handleShelfSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('shelfId').value;
    const data = {
        name: document.getElementById('shelfName').value,
        warehouse_id: document.getElementById('shelfWarehouse').value,
        columns: parseInt(document.getElementById('shelfColumns').value),
        rows: parseInt(document.getElementById('shelfRows').value)
    };
    const url = id ? `${API_BASE}/shelves/${id}` : `${API_BASE}/shelves`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(() => { showNotification(id ? '置物架更新成功' : '置物架添加成功', 'success'); closeModal(); loadShelves(); })
    .catch(() => { showNotification('操作失败', 'error'); });
}

function deleteShelf(id) {
    if (!confirm('确定要删除这个置物架吗？')) return;
    fetch(`${API_BASE}/shelves/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(() => { showNotification('置物架删除成功', 'success'); loadShelves(); })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadBoxes() {
    Promise.all([
        fetch(`${API_BASE}/boxes`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/warehouses`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/shelves`, { headers: getAuthHeader() }).then(r => r.json())
    ])
    .then(([boxes, warehouses, shelves]) => {
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>收纳盒列表</h2>
                    <button class="btn btn-primary" onclick="showBoxModal()">添加收纳盒</button>
                </div>
                <table>
                    <tr><th>编号</th><th>所属空间</th><th>所在置物架</th><th>位置(列-层)</th><th>操作</th></tr>
                    ${boxes.map(b => {
                        const warehouse = warehouses.find(w => w.id === b.warehouse_id);
                        const shelf = shelves.find(s => s.id === b.shelf_id);
                        const pos = b.column > 0 ? b.column + '-' + b.row : '-';
                        return '<tr><td>' + b.box_no + '</td><td>' + (warehouse?.name || '-') + '</td><td>' + (shelf?.name || '-') + '</td><td>' + pos + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="editBox(' + JSON.stringify(b) + ')">编辑</button><button class="btn btn-danger" onclick="deleteBox(\'' + b.id + '\')">删除</button></div></td></tr>';
                    }).join('')}
                </table>
            </div>
        `;
    });
}

function editBox(box) {
    showBoxModal(box);
}

function showBoxModal(box = null) {
    const isEdit = box !== null;
    Promise.all([
        fetch(`${API_BASE}/warehouses`, { headers: getAuthHeader() }).then(r => r.json()),
        fetch(`${API_BASE}/shelves`, { headers: getAuthHeader() }).then(r => r.json())
    ])
    .then(([warehouses, shelves]) => {
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()">
                    <h3>${isEdit ? '编辑收纳盒' : '添加收纳盒'}</h3>
                    <form id="boxForm">
                        <div class="form-group"><label>编号</label><input type="text" id="boxNo" value="${box?.box_no || ''}" required placeholder="如：A001"></div>
                        <div class="form-group"><label>所属储物空间</label><select id="boxWarehouse"><option value="">请选择</option>${warehouses.map(w => '<option value="' + w.id + '"' + (box?.warehouse_id === w.id ? ' selected' : '') + '>' + w.name + '</option>').join('')}</select></div>
                        <div class="form-group"><label>所在置物架</label><select id="boxShelf"><option value="">无</option>${shelves.map(s => '<option value="' + s.id + '"' + (box?.shelf_id === s.id ? ' selected' : '') + '>' + s.name + '</option>').join('')}</select></div>
                        <div class="form-row">
                            <div class="form-group"><label>列</label><input type="number" id="boxColumn" value="${box?.column || ''}" min="1"></div>
                            <div class="form-group"><label>层</label><input type="number" id="boxRow" value="${box?.row || ''}" min="1"></div>
                        </div>
                        <input type="hidden" id="boxId" value="${box?.id || ''}">
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                            <button type="submit" class="btn btn-primary">${isEdit ? '保存' : '添加'}</button>
                        </div>
                    </form>
                </div>
            </div>
        `;
        document.getElementById('boxForm').addEventListener('submit', handleBoxSubmit);
    });
}

function handleBoxSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('boxId').value;
    const data = {
        box_no: document.getElementById('boxNo').value,
        warehouse_id: document.getElementById('boxWarehouse').value,
        shelf_id: document.getElementById('boxShelf').value,
        column: parseInt(document.getElementById('boxColumn').value) || 0,
        row: parseInt(document.getElementById('boxRow').value) || 0
    };
    const url = id ? `${API_BASE}/boxes/${id}` : `${API_BASE}/boxes`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(() => { showNotification(id ? '收纳盒更新成功' : '收纳盒添加成功', 'success'); closeModal(); loadBoxes(); })
    .catch(() => { showNotification('操作失败', 'error'); });
}

function deleteBox(id) {
    if (!confirm('确定要删除这个收纳盒吗？')) return;
    fetch(`${API_BASE}/boxes/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(() => { showNotification('收纳盒删除成功', 'success'); loadBoxes(); })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function logout() {
    localStorage.removeItem('token');
    token = null;
    showLoginPage();
}

function getAuthHeader() {
    return { 'Authorization': 'Bearer ' + token };
}

function showNotification(message, type) {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = 'notification show ' + type;
    setTimeout(() => { notification.classList.remove('show'); }, 3000);
}

function closeModal() {
    document.querySelectorAll('.modal-overlay').forEach(m => m.remove());
}

document.addEventListener('DOMContentLoaded', function() {
    if (token) {
        showMainPage();
    } else {
        showLoginPage();
    }
});