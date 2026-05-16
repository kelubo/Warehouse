const API_BASE = '/api/v1';
let currentPage = { dashboard: 1, products: 1, warehouses: 1, shelves: 1, boxes: 1 };
let totalPages = { products: 1, warehouses: 1, shelves: 1, boxes: 1 };
let searchKeywords = { products: '', warehouses: '', shelves: '', boxes: '' };
let token = localStorage.getItem('access_token') || localStorage.getItem('token');
let refreshToken = localStorage.getItem('refresh_token');
let isRefreshing = false;

function getAuthHeader() {
    return { 'Authorization': 'Bearer ' + token };
}

async function refreshAccessToken() {
    if (isRefreshing) return;
    isRefreshing = true;
    
    try {
        const response = await fetch(`${API_BASE}/refresh-token`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshToken })
        });
        const data = await response.json();
        
        if (data.access_token) {
            localStorage.setItem('access_token', data.access_token);
            token = data.access_token;
        } else {
            logout();
        }
    } catch (error) {
        logout();
    } finally {
        isRefreshing = false;
    }
}

async function fetchWithRefresh(url, options = {}) {
    const response = await fetch(url, {
        ...options,
        headers: {
            ...getAuthHeader(),
            ...options.headers
        }
    });
    
    if (response.status === 401) {
        await refreshAccessToken();
        return fetch(url, {
            ...options,
            headers: {
                ...getAuthHeader(),
                ...options.headers
            }
        });
    }
    
    return response;
}

function showLoginPage() {
    document.getElementById('app').innerHTML = `
        <div class="login-container">
            <div class="login-box">
                <h2>🏠 储物管理</h2>
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
        if (data.access_token) {
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);
            token = data.access_token;
            refreshToken = data.refresh_token;
            showMainPage();
        } else if (data.token) {
            localStorage.setItem('access_token', data.token);
            token = data.token;
            showMainPage();
        } else {
            showNotification('登录失败：' + (data.error || '用户名或密码错误'), 'error');
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
                <div class="logo">🏠 储物管理</div>
                <ul class="menu">
                    <li><a href="#" onclick="loadPage('dashboard')" id="menu-dashboard"><i>📊</i> 总览</a></li>
                    <li><a href="#" onclick="loadPage('products')" id="menu-products"><i>📦</i> 物品管理</a></li>
                    <li><a href="#" onclick="loadPage('warehouses')" id="menu-warehouses"><i>🏠</i> 储物空间</a></li>
                    <li><a href="#" onclick="loadPage('shelves')" id="menu-shelves"><i>📚</i> 置物架</a></li>
                    <li><a href="#" onclick="loadPage('boxes')" id="menu-boxes"><i>📦</i> 收纳盒</a></li>
                    <li><a href="#" onclick="loadPage('units')" id="menu-units"><i>📏</i> 单位管理</a></li>
                    <li><a href="#" onclick="loadPage('categories')" id="menu-categories"><i>🏷️</i> 分类管理</a></li>
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
    document.querySelectorAll('.menu a').forEach(link => link.classList.remove('active'));
    document.getElementById('menu-' + page)?.classList.add('active');
    const titles = {
        dashboard: '总览',
        products: '物品管理',
        warehouses: '储物空间',
        shelves: '置物架',
        boxes: '收纳盒',
        units: '单位管理',
        categories: '分类管理'
    };
    document.getElementById('page-title').textContent = titles[page];
    switch(page) {
        case 'dashboard': loadDashboard(); break;
        case 'products': loadProducts(); break;
        case 'warehouses': loadWarehouses(); break;
        case 'shelves': loadShelves(); break;
        case 'boxes': loadBoxes(); break;
        case 'units': loadUnits(); break;
        case 'categories': loadCategories(); break;
    }
}

function loadDashboard() {
    Promise.all([
        fetch(`${API_BASE}/products?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('products'); return r.json(); }),
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('warehouses'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('shelves'); return r.json(); }),
        fetch(`${API_BASE}/boxes?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('boxes'); return r.json(); })
    ])
    .then(([productsResult, warehousesResult, shelvesResult, boxesResult]) => {
        const products = productsResult.items || [];
        const warehouses = warehousesResult.items || [];
        const shelves = shelvesResult.items || [];
        const boxes = boxesResult.items || [];

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
                    ${products.length === 0 ? '<tr><td colspan="3" style="text-align:center">暂无物品，请添加物品</td></tr>' : products.slice(-5).reverse().map(p => {
                        const shelf = shelves.find(s => s.id === p.shelf_id);
                        const box = boxes.find(b => b.id === p.box_id);
                        const warehouse = warehouses.find(w => w.id === (box?.warehouse_id || shelf?.warehouse_id));
                        let location = '未指定';
                        if (box && shelf) {
                            location = (warehouse?.name || '') + ' > ' + shelf.name + ' (列' + box.column + '-层' + box.row + ')' + ' > 收纳盒:' + box.box_no;
                        } else if (shelf) {
                            location = (warehouse?.name || '') + ' > ' + shelf.name;
                            if (p.shelf_column > 0) location += ' (列' + p.shelf_column + '-层' + p.shelf_row + ')';
                        }
                        return `<tr><td>${p.name}</td><td>${p.category || '-'}</td><td>${location}</td></tr>`;
                    }).join('')}
                </table>
            </div>
        `;
    })
    .catch(error => {
        console.error('Dashboard load error:', error);
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <h2>加载失败</h2>
                <p>请检查是否已登录或联系管理员</p>
                <button class="btn btn-primary" onclick="localStorage.removeItem('token'); localStorage.removeItem('access_token'); location.reload()">重新登录</button>
            </div>
        `;
    });
}

function renderPagination(type, total, current) {
    if (total <= 1) return '';
    let html = '<div class="pagination">';
    html += `<button class="btn btn-secondary" ${current <= 1 ? 'disabled' : ''} onclick="load${type.charAt(0).toUpperCase() + type.slice(1)}(${current - 1})">上一页</button>`;
    for (let i = 1; i <= total; i++) {
        html += `<button class="btn ${i === current ? 'btn-primary' : 'btn-secondary'}" onclick="load${type.charAt(0).toUpperCase() + type.slice(1)}(${i})">${i}</button>`;
    }
    html += `<button class="btn btn-secondary" ${current >= total ? 'disabled' : ''} onclick="load${type.charAt(0).toUpperCase() + type.slice(1)}(${current + 1})">下一页</button>`;
    html += '</div>';
    return html;
}

function loadProducts(page = 1, keyword = '') {
    currentPage.products = page;
    if (keyword !== undefined) searchKeywords.products = keyword;

    Promise.all([
        fetch(`${API_BASE}/products?page=${page}&page_size=10&keyword=${encodeURIComponent(searchKeywords.products)}`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取物品列表失败'); return r.json(); }),
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); }),
        fetch(`${API_BASE}/boxes?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取收纳盒列表失败'); return r.json(); })
    ])
    .then(([result, warehousesResult, shelvesResult, boxesResult]) => {
        const products = result.items || result;
        const warehouses = warehousesResult.items || warehousesResult;
        const shelves = shelvesResult.items || shelvesResult;
        const boxes = boxesResult.items || boxesResult;
        totalPages.products = result.total_pages || 1;

        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>物品列表</h2>
                    <div class="search-bar">
                        <input type="text" id="productSearch" placeholder="搜索物品名称、SKU或分类..." value="${searchKeywords.products}">
                        <button class="btn btn-secondary" onclick="loadProducts(1, document.getElementById('productSearch').value)">搜索</button>
                        <button class="btn btn-primary" onclick="showProductModal()">添加物品</button>
                    </div>
                </div>
                <table>
                    <tr><th>名称</th><th>分类</th><th>数量</th><th>单位</th><th>存放位置</th><th>操作</th></tr>
                    ${products.length === 0 ? '<tr><td colspan="6" style="text-align:center;">暂无数据</td></tr>' : products.map(p => {
                        const shelf = shelves.find(s => s.id === p.shelf_id);
                        const box = boxes.find(b => b.id === p.box_id);
                        const warehouse = warehouses.find(w => w.id === (box?.warehouse_id || shelf?.warehouse_id));
                        let location = '未指定';
                        if (box && shelf) {
                            location = (warehouse?.name || '') + ' > ' + shelf.name + ' (列' + box.column + '-层' + box.row + ')' + ' > 收纳盒:' + box.box_no;
                        } else if (shelf) {
                            location = (warehouse?.name || '') + ' > ' + shelf.name;
                            if (p.shelf_column > 0) location += ' (列' + p.shelf_column + '-层' + p.shelf_row + ')';
                        }
                        return '<tr><td>' + p.name + '</td><td>' + (p.category || '-') + '</td><td>' + (p.quantity || 1) + '</td><td>' + (p.unit || '-') + '</td><td>' + location + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="showProductDetail(\'' + p.id + '\')">详情</button><button class="btn btn-secondary" onclick="editProduct(\'' + p.id + '\')">编辑</button><button class="btn btn-danger" onclick="deleteProduct(\'' + p.id + '\')">删除</button></div></td></tr>';
                    }).join('')}
                </table>
                ${renderPagination('products', totalPages.products, currentPage.products)}
            </div>
        `;
    })
    .catch(error => {
        showNotification('加载物品列表失败: ' + error.message, 'error');
    });
}

function editProduct(id) {
    fetch(`${API_BASE}/products/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(product => {
        showProductModal(product);
    })
    .catch(error => {
        showNotification('获取物品信息失败', 'error');
    });
}

function showProductDetail(id) {
    Promise.all([
        fetch(`${API_BASE}/products/${id}`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取物品信息失败'); return r.json(); }),
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); }),
        fetch(`${API_BASE}/boxes?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取收纳盒列表失败'); return r.json(); })
    ])
    .then(([product, warehousesResult, shelvesResult, boxesResult]) => {
        const warehouses = warehousesResult.items || warehousesResult;
        const shelves = shelvesResult.items || shelvesResult;
        const boxes = boxesResult.items || boxesResult;
        const warehouse = warehouses.find(w => w.id === product.warehouse_id);
        const shelf = shelves.find(s => s.id === product.shelf_id);
        const box = boxes.find(b => b.id === product.box_id);

        let location = '-';
        if (box) {
            location = box.box_no;
        } else if (shelf) {
            location = shelf.name;
            if (product.shelf_column > 0) location += ' (列' + product.shelf_column + '-层' + product.shelf_row + ')';
        } else if (warehouse) {
            location = warehouse.name;
        }

        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
                    <h3>物品详情</h3>
                    <div style="margin-bottom: 15px;">
                        <p><strong>名称：</strong>${product.name}</p>
                        <p><strong>分类：</strong>${product.category || '-'}</p>
                        <p><strong>数量：</strong>${product.quantity || 1}</p>
                        <p><strong>存放位置：</strong>${location}</p>
                        <p><strong>描述：</strong>${product.description || '-'}</p>
                    </div>
                    <div class="modal-footer" style="margin-top: 20px;">
                        <button type="button" class="btn btn-secondary" onclick="closeModal()">关闭</button>
                    </div>
                </div>
            </div>
        `;
    })
    .catch(error => {
        showNotification('获取物品详情失败', 'error');
    });
}

function updateShelfPositionVisibility() {
    const shelfSelect = document.getElementById('productShelf');
    const boxSelect = document.getElementById('productBox');
    const columnInput = document.getElementById('productShelfColumn');
    const rowInput = document.getElementById('productShelfRow');
    const shelfPositionDiv = document.getElementById('shelfPosition');
    
    const shelfId = shelfSelect.value;
    const boxId = boxSelect.value;
    
    if (shelfId && !boxId) {
        shelfPositionDiv.style.display = '';
        // 设置行列限制
        const selectedOption = shelfSelect.options[shelfSelect.selectedIndex];
        const maxColumns = parseInt(selectedOption.getAttribute('data-columns')) || 999;
        const maxRows = parseInt(selectedOption.getAttribute('data-rows')) || 999;
        
        columnInput.max = maxColumns;
        rowInput.max = maxRows;
        
        if (columnInput.value && parseInt(columnInput.value) > maxColumns) {
            columnInput.value = maxColumns;
        }
        if (rowInput.value && parseInt(rowInput.value) > maxRows) {
            rowInput.value = maxRows;
        }
    } else {
        shelfPositionDiv.style.display = 'none';
    }
}

function showProductModal(product = null) {
    const isEdit = product !== null;
    Promise.all([
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); }),
        fetch(`${API_BASE}/boxes?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取收纳盒列表失败'); return r.json(); }),
        fetch(`${API_BASE}/units`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取单位列表失败'); return r.json(); }),
        fetch(`${API_BASE}/categories`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取分类列表失败'); return r.json(); })
    ])
    .then(([warehousesResult, shelvesResult, boxesResult, unitsResult, categoriesResult]) => {
        const warehouses = warehousesResult.items || warehousesResult;
        const shelves = shelvesResult.items || shelvesResult;
        const boxes = boxesResult.items || boxesResult;
        const units = unitsResult || [];
        const categories = categoriesResult || [];
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()">
                    <h3>${isEdit ? '编辑物品' : '添加物品'}</h3>
                    <form id="productForm">
                        <div class="form-group"><label>物品名称</label><input type="text" id="productName" value="${product?.name || ''}" required></div>
                        <div class="form-row">
                            <div class="form-group"><label>分类</label><select id="productCategory"><option value="">请选择</option>${categories.map(c => '<option value="' + c.name + '"' + (product?.category === c.name ? ' selected' : '') + '>' + c.name + '</option>').join('')}</select></div>
                            <div class="form-group"><label>数量</label><input type="number" id="productQuantity" value="${product?.quantity || 1}" min="1" style="width:80px;"></div>
                            <div class="form-group"><label>单位</label><select id="productUnit" style="width:80px;"><option value="">请选择</option>${units.map(u => '<option value="' + u.name + '"' + (product?.unit === u.name ? ' selected' : '') + '>' + u.name + '</option>').join('')}</select></div>
                        </div>
                        <div class="form-group"><label>存放位置 - 储物空间</label><select id="productWarehouse"><option value="">无</option>${warehouses.map(w => '<option value="' + w.id + '"' + (product?.warehouse_id === w.id ? ' selected' : '') + '>' + w.name + '</option>').join('')}</select></div>
                        <div class="form-group"><label>存放位置 - 置物架</label><select id="productShelf"><option value="">无</option>${shelves.map(s => '<option value="' + s.id + '" data-columns="' + s.columns + '" data-rows="' + s.rows + '"' + (product?.shelf_id === s.id ? ' selected' : '') + '>' + s.name + '</option>').join('')}</select></div>
                        <div class="form-group"><label>存放位置 - 收纳盒</label><select id="productBox"><option value="">无</option>${boxes.map(b => '<option value="' + b.id + '"' + (product?.box_id === b.id ? ' selected' : '') + '>' + b.box_no + '</option>').join('')}</select></div>
                        <div class="form-row" id="shelfPosition" ${(product?.shelf_id && !product?.box_id) ? '' : 'style="display:none"'}>
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
        document.getElementById('productShelf').addEventListener('change', updateShelfPositionVisibility);
        document.getElementById('productBox').addEventListener('change', updateShelfPositionVisibility);
    })
    .catch(error => {
        showNotification('加载表单数据失败: ' + error.message, 'error');
    });
}

function handleProductSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('productId').value;
    const data = {
        name: document.getElementById('productName').value,
        category: document.getElementById('productCategory').value,
        unit: document.getElementById('productUnit').value,
        quantity: parseInt(document.getElementById('productQuantity').value) || 1,
        warehouse_id: document.getElementById('productWarehouse').value,
        shelf_id: document.getElementById('productShelf').value,
        box_id: document.getElementById('productBox').value,
        shelf_column: parseInt(document.getElementById('productShelfColumn').value) || 0,
        shelf_row: parseInt(document.getElementById('productShelfRow').value) || 0,
        description: document.getElementById('productDescription').value
    };
    const url = id ? `${API_BASE}/products/${id}` : `${API_BASE}/products`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(response => {
        if (response.ok) {
            showNotification(id ? '物品更新成功' : '物品添加成功', 'success');
            closeModal();
            loadProducts();
        } else {
            response.json().then(data => {
                showNotification(data.error || '操作失败', 'error');
            }).catch(() => {
                showNotification('操作失败', 'error');
            });
        }
    })
    .catch(() => { showNotification('操作失败', 'error'); });
}

function deleteProduct(id) {
    if (!confirm('确定要删除这个物品吗？')) return;
    fetch(`${API_BASE}/products/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(response => {
        return response.json().then(data => ({ response, data }));
    })
    .then(({ response, data }) => {
        if (response.ok) {
            showNotification('物品删除成功', 'success');
            loadProducts();
        } else {
            showNotification(data.error || '删除失败', 'error');
        }
    })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadWarehouses(page = 1, keyword = '') {
    currentPage.warehouses = page;
    if (keyword !== undefined) searchKeywords.warehouses = keyword;

    fetch(`${API_BASE}/warehouses?page=${page}&page_size=10&keyword=${encodeURIComponent(searchKeywords.warehouses)}`, { headers: getAuthHeader() })
    .then(response => {
        if (!response.ok) {
            throw new Error('获取储物空间列表失败');
        }
        return response.json();
    })
    .then(result => {
        const warehouses = result.items || result;
        totalPages.warehouses = result.total_pages || 1;

        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>储物空间列表</h2>
                    <div class="action-bar">
                        <button class="btn btn-primary" onclick="showWarehouseModal()">添加储物空间</button>
                    </div>
                </div>
                <table>
                    <tr><th>名称</th><th>位置</th><th>容量</th><th>描述</th><th>操作</th></tr>
                    ${warehouses.length === 0 ? '<tr><td colspan="5" style="text-align:center;">暂无数据</td></tr>' : warehouses.map(w => '<tr><td>' + w.name + '</td><td>' + (w.location || '-') + '</td><td>' + (w.capacity || '-') + '</td><td>' + (w.description || '-') + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="showWarehouseDetail(\'' + w.id + '\')">详情</button><button class="btn btn-secondary" onclick="editWarehouse(\'' + w.id + '\')">编辑</button><button class="btn btn-danger" onclick="deleteWarehouse(\'' + w.id + '\')">删除</button></div></td></tr>').join('')}
                </table>
                ${renderPagination('warehouses', totalPages.warehouses, currentPage.warehouses)}
            </div>
        `;
    })
    .catch(error => {
        showNotification('加载储物空间列表失败: ' + error.message, 'error');
    });
}

function editWarehouse(id) {
    fetch(`${API_BASE}/warehouses/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(warehouse => {
        showWarehouseModal(warehouse);
    })
    .catch(error => {
        showNotification('获取储物空间信息失败', 'error');
    });
}

function showWarehouseDetail(id) {
    fetch(`${API_BASE}/warehouses/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(warehouse => {
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
                    <h3>储物空间详情</h3>
                    <div style="margin-bottom: 15px;">
                        <p><strong>名称：</strong>${warehouse.name}</p>
                        <p><strong>位置：</strong>${warehouse.location || '-'}</p>
                        <p><strong>容量：</strong>${warehouse.capacity || '-'}</p>
                        <p><strong>描述：</strong>${warehouse.description || '-'}</p>
                    </div>
                    ${warehouse.shelves && warehouse.shelves.length > 0 ? `
                    <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                        <h4>包含的置物架 (${warehouse.shelves.length})</h4>
                        <table style="width:100%; margin-top:10px;">
                            <tr><th>名称</th><th>规格(列×层)</th></tr>
                            ${warehouse.shelves.map(s => '<tr><td>' + s.name + '</td><td>' + s.columns + ' × ' + s.rows + '</td></tr>').join('')}
                        </table>
                    </div>
                    ` : '<p style="color:#999; margin-top:15px;">暂无置物架</p>'}
                    <div class="modal-footer" style="margin-top: 20px;">
                        <button type="button" class="btn btn-secondary" onclick="closeModal()">关闭</button>
                    </div>
                </div>
            </div>
        `;
    })
    .catch(error => {
        showNotification('获取储物空间详情失败', 'error');
    });
}

function showWarehouseModal(warehouse = null) {
    const isEdit = warehouse !== null;
    document.getElementById('page-content').innerHTML += `
        <div class="modal-overlay active" onclick="closeModal()">
            <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
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
                ${isEdit && warehouse?.shelves && warehouse.shelves.length > 0 ? `
                <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                    <h4>包含的置物架 (${warehouse.shelves.length})</h4>
                    <table style="width:100%; margin-top:10px;">
                        <tr><th>名称</th><th>规格(列×层)</th></tr>
                        ${warehouse.shelves.map(s => '<tr><td>' + s.name + '</td><td>' + s.columns + ' × ' + s.rows + '</td></tr>').join('')}
                    </table>
                </div>
                ` : ''}
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
    .then(response => {
        return response.json().then(data => ({ response, data }));
    })
    .then(({ response, data }) => {
        if (response.ok) {
            showNotification('储物空间删除成功', 'success');
            loadWarehouses();
        } else {
            showNotification(data.error || '删除失败', 'error');
        }
    })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadShelves(page = 1, keyword = '') {
    currentPage.shelves = page;
    if (keyword !== undefined) searchKeywords.shelves = keyword;

    Promise.all([
        fetch(`${API_BASE}/shelves?page=${page}&page_size=10&keyword=${encodeURIComponent(searchKeywords.shelves)}`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); }),
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); })
    ])
    .then(([result, warehousesResult]) => {
        const shelves = result.items || result;
        const warehouses = warehousesResult.items || warehousesResult;
        totalPages.shelves = result.total_pages || 1;

        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>置物架列表</h2>
                    <div class="action-bar">
                        <button class="btn btn-primary" onclick="showShelfModal()">添加置物架</button>
                    </div>
                </div>
                <table>
                    <tr><th>名称</th><th>所属空间</th><th>列数</th><th>层数</th><th>操作</th></tr>
                    ${shelves.length === 0 ? '<tr><td colspan="5" style="text-align:center;">暂无数据</td></tr>' : shelves.map(s => {
                        const warehouse = warehouses.find(w => w.id === s.warehouse_id);
                        return '<tr><td>' + s.name + '</td><td>' + (warehouse?.name || '-') + '</td><td>' + s.columns + '</td><td>' + s.rows + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="showShelfDetail(\'' + s.id + '\')">详情</button><button class="btn btn-secondary" onclick="editShelf(\'' + s.id + '\')">编辑</button><button class="btn btn-danger" onclick="deleteShelf(\'' + s.id + '\')">删除</button></div></td></tr>';
                    }).join('')}
                </table>
                ${renderPagination('shelves', totalPages.shelves, currentPage.shelves)}
            </div>
        `;
    })
    .catch(error => {
        showNotification('加载置物架列表失败: ' + error.message, 'error');
    });
}

function editShelf(id) {
    fetch(`${API_BASE}/shelves/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(shelf => {
        showShelfModal(shelf);
    })
    .catch(error => {
        showNotification('获取置物架信息失败', 'error');
    });
}

function showShelfDetail(id) {
    fetch(`${API_BASE}/shelves/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(shelf => {
        fetch(`${API_BASE}/warehouses/${shelf.warehouse_id}`, { headers: getAuthHeader() })
        .then(response => response.json())
        .then(warehouse => {
            document.getElementById('page-content').innerHTML += `
                <div class="modal-overlay active" onclick="closeModal()">
                    <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
                        <h3>置物架详情</h3>
                        <div style="margin-bottom: 15px;">
                            <p><strong>名称：</strong>${shelf.name}</p>
                            <p><strong>所属储物空间：</strong>${warehouse.name}</p>
                            <p><strong>规格：</strong>${shelf.columns} 列 × ${shelf.rows} 层</p>
                        </div>
                        ${shelf.boxes && shelf.boxes.length > 0 ? `
                        <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                            <h4>包含的收纳盒 (${shelf.boxes.length})</h4>
                            <table style="width:100%; margin-top:10px;">
                                <tr><th>编号</th><th>位置(列-层)</th></tr>
                                ${shelf.boxes.map(b => '<tr><td>' + b.box_no + '</td><td>' + (b.column > 0 ? b.column + '-' + b.row : '-') + '</td></tr>').join('')}
                            </table>
                        </div>
                        ` : '<p style="color:#999; margin-top:15px;">暂无收纳盒</p>'}
                        <div class="modal-footer" style="margin-top: 20px;">
                            <button type="button" class="btn btn-secondary" onclick="closeModal()">关闭</button>
                        </div>
                    </div>
                </div>
            `;
        })
        .catch(() => {
            document.getElementById('page-content').innerHTML += `
                <div class="modal-overlay active" onclick="closeModal()">
                    <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
                        <h3>置物架详情</h3>
                        <div style="margin-bottom: 15px;">
                            <p><strong>名称：</strong>${shelf.name}</p>
                            <p><strong>所属储物空间：</strong>-</p>
                            <p><strong>规格：</strong>${shelf.columns} 列 × ${shelf.rows} 层</p>
                        </div>
                        ${shelf.boxes && shelf.boxes.length > 0 ? `
                        <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                            <h4>包含的收纳盒 (${shelf.boxes.length})</h4>
                            <table style="width:100%; margin-top:10px;">
                                <tr><th>编号</th><th>位置(列-层)</th></tr>
                                ${shelf.boxes.map(b => '<tr><td>' + b.box_no + '</td><td>' + (b.column > 0 ? b.column + '-' + b.row : '-') + '</td></tr>').join('')}
                            </table>
                        </div>
                        ` : '<p style="color:#999; margin-top:15px;">暂无收纳盒</p>'}
                        <div class="modal-footer" style="margin-top: 20px;">
                            <button type="button" class="btn btn-secondary" onclick="closeModal()">关闭</button>
                        </div>
                    </div>
                </div>
            `;
        });
    })
    .catch(error => {
        showNotification('获取置物架详情失败', 'error');
    });
}

function showShelfModal(shelf = null) {
    const isEdit = shelf !== null;
    fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() })
    .then(response => { if (!response.ok) throw new Error('获取储物空间列表失败'); return response.json(); })
    .then(warehousesResult => {
        const warehouses = warehousesResult.items || warehousesResult;
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
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
                    ${isEdit && shelf?.boxes && shelf.boxes.length > 0 ? `
                    <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                        <h4>包含的收纳盒 (${shelf.boxes.length})</h4>
                        <table style="width:100%; margin-top:10px;">
                            <tr><th>编号</th><th>位置(列-层)</th></tr>
                            ${shelf.boxes.map(b => '<tr><td>' + b.box_no + '</td><td>' + (b.column > 0 ? b.column + '-' + b.row : '-') + '</td></tr>').join('')}
                        </table>
                    </div>
                    ` : ''}
                </div>
            </div>
        `;
        document.getElementById('shelfForm').addEventListener('submit', handleShelfSubmit);
    })
    .catch(error => {
        showNotification('加载表单数据失败: ' + error.message, 'error');
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
    .then(response => {
        return response.json().then(data => ({ response, data }));
    })
    .then(({ response, data }) => {
        if (response.ok) {
            showNotification('置物架删除成功', 'success');
            loadShelves();
        } else {
            showNotification(data.error || '删除失败', 'error');
        }
    })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadUnits() {
    fetch(`${API_BASE}/units`, { headers: getAuthHeader() })
    .then(r => { if (!r.ok) throw new Error('获取单位列表失败'); return r.json(); })
    .then(units => {
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>单位列表</h2>
                    <div class="action-bar">
                        <button class="btn btn-primary" onclick="showUnitModal()">添加单位</button>
                    </div>
                </div>
                <table>
                    <tr><th>名称</th><th>操作</th></tr>
                    ${units.length === 0 ? '<tr><td colspan="2" style="text-align:center;">暂无数据</td></tr>' : units.map(u => '<tr><td>' + u.name + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="editUnit(\'' + u.id + '\')">编辑</button><button class="btn btn-danger" onclick="deleteUnit(\'' + u.id + '\')">删除</button></div></td></tr>').join('')}
                </table>
            </div>
        `;
    })
    .catch(error => {
        showNotification('加载单位列表失败: ' + error.message, 'error');
    });
}

function showUnitModal(unit = null) {
    const isEdit = unit !== null;
    document.getElementById('page-content').innerHTML += `
        <div class="modal-overlay active" onclick="closeModal()">
            <div class="modal" onclick="event.stopPropagation()">
                <h3>${isEdit ? '编辑单位' : '添加单位'}</h3>
                <form id="unitForm">
                    <div class="form-group"><label>名称</label><input type="text" id="unitName" value="${unit?.name || ''}" required placeholder="如：个、箱、件、千克"></div>
                    <input type="hidden" id="unitId" value="${unit?.id || ''}">
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                        <button type="submit" class="btn btn-primary">${isEdit ? '保存' : '添加'}</button>
                    </div>
                </form>
            </div>
        </div>
    `;
    document.getElementById('unitForm').addEventListener('submit', handleUnitSubmit);
}

function handleUnitSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('unitId').value;
    const data = { name: document.getElementById('unitName').value };
    const url = id ? `${API_BASE}/units/${id}` : `${API_BASE}/units`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(response => {
        if (response.ok) {
            showNotification(id ? '单位更新成功' : '单位添加成功', 'success');
            closeModal();
            loadUnits();
        } else {
            return response.json().then(data => { throw new Error(data.error || '操作失败'); });
        }
    })
    .catch(error => { showNotification(error.message, 'error'); });
}

function editUnit(id) {
    fetch(`${API_BASE}/units/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(unit => { showUnitModal(unit); })
    .catch(() => { showNotification('获取单位信息失败', 'error'); });
}

function deleteUnit(id) {
    if (!confirm('确定要删除这个单位吗？')) return;
    fetch(`${API_BASE}/units/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(response => {
        return response.json().then(data => ({ response, data }));
    })
    .then(({ response, data }) => {
        if (response.ok) {
            showNotification('单位删除成功', 'success');
            loadUnits();
        } else {
            showNotification(data.error || '删除失败', 'error');
        }
    })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadCategories() {
    fetch(`${API_BASE}/categories`, { headers: getAuthHeader() })
    .then(r => { if (!r.ok) throw new Error('获取分类列表失败'); return r.json(); })
    .then(categories => {
        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>分类列表</h2>
                    <div class="action-bar">
                        <button class="btn btn-primary" onclick="showCategoryModal()">添加分类</button>
                    </div>
                </div>
                <table>
                    <tr><th>名称</th><th>操作</th></tr>
                    ${categories.length === 0 ? '<tr><td colspan="2" style="text-align:center;">暂无数据</td></tr>' : categories.map(c => '<tr><td>' + c.name + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="editCategory(\'' + c.id + '\')">编辑</button><button class="btn btn-danger" onclick="deleteCategory(\'' + c.id + '\')">删除</button></div></td></tr>').join('')}
                </table>
            </div>
        `;
    })
    .catch(error => {
        showNotification('加载分类列表失败: ' + error.message, 'error');
    });
}

function showCategoryModal(category = null) {
    const isEdit = category !== null;
    document.getElementById('page-content').innerHTML += `
        <div class="modal-overlay active" onclick="closeModal()">
            <div class="modal" onclick="event.stopPropagation()">
                <h3>${isEdit ? '编辑分类' : '添加分类'}</h3>
                <form id="categoryForm">
                    <div class="form-group"><label>名称</label><input type="text" id="categoryName" value="${category?.name || ''}" required placeholder="如：电子产品、服装、食品"></div>
                    <input type="hidden" id="categoryId" value="${category?.id || ''}">
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                        <button type="submit" class="btn btn-primary">${isEdit ? '保存' : '添加'}</button>
                    </div>
                </form>
            </div>
        </div>
    `;
    document.getElementById('categoryForm').addEventListener('submit', handleCategorySubmit);
}

function handleCategorySubmit(e) {
    e.preventDefault();
    const id = document.getElementById('categoryId').value;
    const data = { name: document.getElementById('categoryName').value };
    const url = id ? `${API_BASE}/categories/${id}` : `${API_BASE}/categories`;
    const method = id ? 'PUT' : 'POST';
    fetch(url, { method, headers: { ...getAuthHeader(), 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
    .then(response => {
        if (response.ok) {
            showNotification(id ? '分类更新成功' : '分类添加成功', 'success');
            closeModal();
            loadCategories();
        } else {
            return response.json().then(data => { throw new Error(data.error || '操作失败'); });
        }
    })
    .catch(error => { showNotification(error.message, 'error'); });
}

function editCategory(id) {
    fetch(`${API_BASE}/categories/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(category => { showCategoryModal(category); })
    .catch(() => { showNotification('获取分类信息失败', 'error'); });
}

function deleteCategory(id) {
    if (!confirm('确定要删除这个分类吗？')) return;
    fetch(`${API_BASE}/categories/${id}`, { method: 'DELETE', headers: getAuthHeader() })
    .then(response => {
        return response.json().then(data => ({ response, data }));
    })
    .then(({ response, data }) => {
        if (response.ok) {
            showNotification('分类删除成功', 'success');
            loadCategories();
        } else {
            showNotification(data.error || '删除失败', 'error');
        }
    })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function loadBoxes(page = 1) {
    currentPage.boxes = page;

    Promise.all([
        fetch(`${API_BASE}/boxes?page=${page}&page_size=10`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取收纳盒列表失败'); return r.json(); }),
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); })
    ])
    .then(([result, warehousesResult, shelvesResult]) => {
        const boxes = result.items || result;
        const warehouses = warehousesResult.items || warehousesResult;
        const shelves = shelvesResult.items || shelvesResult;
        totalPages.boxes = result.total_pages || 1;

        document.getElementById('page-content').innerHTML = `
            <div class="card">
                <div class="card-header">
                    <h2>收纳盒列表</h2>
                    <div class="action-bar">
                        <button class="btn btn-primary" onclick="showBoxModal()">添加收纳盒</button>
                    </div>
                </div>
                <table>
                    <tr><th>编号</th><th>所属空间</th><th>所在置物架</th><th>位置(列-层)</th><th>操作</th></tr>
                    ${boxes.length === 0 ? '<tr><td colspan="5" style="text-align:center;">暂无数据</td></tr>' : boxes.map(b => {
                        const warehouse = warehouses.find(w => w.id === b.warehouse_id);
                        const shelf = shelves.find(s => s.id === b.shelf_id);
                        const pos = b.column > 0 ? b.column + '-' + b.row : '-';
                        return '<tr><td>' + b.box_no + '</td><td>' + (warehouse?.name || '-') + '</td><td>' + (shelf?.name || '-') + '</td><td>' + pos + '</td><td><div class="btn-group"><button class="btn btn-secondary" onclick="showBoxDetail(\'' + b.id + '\')">详情</button><button class="btn btn-secondary" onclick="editBox(\'' + b.id + '\')">编辑</button><button class="btn btn-danger" onclick="deleteBox(\'' + b.id + '\')">删除</button></div></td></tr>';
                    }).join('')}
                </table>
                ${renderPagination('boxes', totalPages.boxes, currentPage.boxes)}
            </div>
        `;
    })
    .catch(error => {
        showNotification('加载收纳盒列表失败: ' + error.message, 'error');
    });
}

function editBox(id) {
    fetch(`${API_BASE}/boxes/${id}`, { headers: getAuthHeader() })
    .then(response => response.json())
    .then(box => {
        showBoxModal(box);
    })
    .catch(error => {
        showNotification('获取收纳盒信息失败', 'error');
    });
}

function showBoxDetail(id) {
    Promise.all([
        fetch(`${API_BASE}/boxes/${id}`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取收纳盒信息失败'); return r.json(); }),
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); })
    ])
    .then(([box, warehousesResult, shelvesResult]) => {
        const warehouses = warehousesResult.items || warehousesResult;
        const shelves = shelvesResult.items || shelvesResult;
        const warehouse = warehouses.find(w => w.id === box.warehouse_id);
        const shelf = shelves.find(s => s.id === box.shelf_id);
        const pos = box.column > 0 ? box.column + '-' + box.row : '-';

        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
                    <h3>收纳盒详情</h3>
                    <div style="margin-bottom: 15px;">
                        <p><strong>编号：</strong>${box.box_no}</p>
                        <p><strong>所属储物空间：</strong>${warehouse?.name || '-'}</p>
                        <p><strong>所在置物架：</strong>${shelf?.name || '-'}</p>
                        <p><strong>位置(列-层)：</strong>${pos}</p>
                    </div>
                    ${box.products && box.products.length > 0 ? `
                    <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                        <h4>包含的物品 (${box.products.length})</h4>
                        <table style="width:100%; margin-top:10px;">
                            <tr><th>名称</th><th>分类</th><th>数量</th><th>单位</th></tr>
                            ${box.products.map(p => '<tr><td>' + p.name + '</td><td>' + (p.category || '-') + '</td><td>' + p.quantity + '</td><td>' + (p.unit || '-') + '</td></tr>').join('')}
                        </table>
                    </div>
                    ` : '<p style="color:#999; margin-top:15px;">暂无物品</p>'}
                    <div class="modal-footer" style="margin-top: 20px;">
                        <button type="button" class="btn btn-secondary" onclick="closeModal()">关闭</button>
                    </div>
                </div>
            </div>
        `;
    })
    .catch(error => {
        showNotification('获取收纳盒详情失败', 'error');
    });
}

function showBoxModal(box = null) {
    const isEdit = box !== null;
    Promise.all([
        fetch(`${API_BASE}/warehouses?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取储物空间列表失败'); return r.json(); }),
        fetch(`${API_BASE}/shelves?page=1&page_size=1000`, { headers: getAuthHeader() }).then(r => { if (!r.ok) throw new Error('获取置物架列表失败'); return r.json(); })
    ])
    .then(([warehousesResult, shelvesResult]) => {
        const warehouses = warehousesResult.items || warehousesResult;
        const shelves = shelvesResult.items || shelvesResult;
        document.getElementById('page-content').innerHTML += `
            <div class="modal-overlay active" onclick="closeModal()">
                <div class="modal" onclick="event.stopPropagation()" style="max-width: 800px;">
                    <h3>${isEdit ? '编辑收纳盒' : '添加收纳盒'}</h3>
                    <form id="boxForm">
                        <div class="form-group"><label>编号</label><input type="text" id="boxNo" value="${box?.box_no || ''}" required placeholder="如：A001"></div>
                        <div class="form-group"><label>所属储物空间</label><select id="boxWarehouse"><option value="">请选择</option>${warehouses.map(w => '<option value="' + w.id + '"' + (box?.warehouse_id === w.id ? ' selected' : '') + '>' + w.name + '</option>').join('')}</select></div>
                        <div class="form-group"><label>所在置物架</label><select id="boxShelf"><option value="">无</option>${shelves.map(s => '<option value="' + s.id + '" data-columns="' + s.columns + '" data-rows="' + s.rows + '"' + (box?.shelf_id === s.id ? ' selected' : '') + '>' + s.name + '</option>').join('')}</select></div>
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
                    ${isEdit && box?.products && box.products.length > 0 ? `
                    <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
                        <h4>包含的物品 (${box.products.length})</h4>
                        <table style="width:100%; margin-top:10px;">
                            <tr><th>名称</th><th>分类</th><th>数量</th><th>单位</th></tr>
                            ${box.products.map(p => '<tr><td>' + p.name + '</td><td>' + (p.category || '-') + '</td><td>' + p.quantity + '</td><td>' + (p.unit || '-') + '</td></tr>').join('')}
                        </table>
                    </div>
                    ` : ''}
                </div>
            </div>
        `;
        document.getElementById('boxForm').addEventListener('submit', handleBoxSubmit);

        const shelfSelect = document.getElementById('boxShelf');
        const columnInput = document.getElementById('boxColumn');
        const rowInput = document.getElementById('boxRow');

        function updateMaxValues() {
            const selectedOption = shelfSelect.options[shelfSelect.selectedIndex];
            const maxColumns = parseInt(selectedOption.getAttribute('data-columns')) || 999;
            const maxRows = parseInt(selectedOption.getAttribute('data-rows')) || 999;

            columnInput.max = maxColumns;
            rowInput.max = maxRows;

            if (columnInput.value && parseInt(columnInput.value) > maxColumns) {
                columnInput.value = maxColumns;
            }
            if (rowInput.value && parseInt(rowInput.value) > maxRows) {
                rowInput.value = maxRows;
            }
        }

        updateMaxValues();
        shelfSelect.addEventListener('change', updateMaxValues);
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
    .then(response => {
        return response.json().then(data => ({ response, data }));
    })
    .then(({ response, data }) => {
        if (response.ok) {
            showNotification('收纳盒删除成功', 'success');
            loadBoxes();
        } else {
            showNotification(data.error || '删除失败', 'error');
        }
    })
    .catch(() => { showNotification('删除失败', 'error'); });
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    token = null;
    refreshToken = null;
    showLoginPage();
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