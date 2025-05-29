// ===========================================
// ADMIN DASHBOARD - SIMPLIFIED VERSION
// ===========================================

// Global variables
let currentTab = 'pending';
const API_BASE = 'http://localhost:8080/api';

// ===========================================
// LOGIN FUNCTIONALITY
// ===========================================
if (window.location.pathname.endsWith('login.html')) {
    const loginForm = document.getElementById("adminLoginForm");
    loginForm.addEventListener("submit", function(e) {
        e.preventDefault();

        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;

        fetch(`${API_BASE}/admin/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: 'include',
            body: JSON.stringify({
                username,
                password,
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Gagal login")
            }
            return response.json();
        })
        .then(data => {
            console.log("Login berhasil", data);
            sessionStorage.setItem("user-admin", JSON.stringify(data));
            window.location.href = "index.html";
        })
        .catch(error => {
            console.log(error);
            alert('Login gagal: ' + error.message);
        });
    });
}

// ===========================================
// DASHBOARD FUNCTIONALITY
// ===========================================
if (window.location.pathname.endsWith('index.html')) {
    
    // Check authentication on load
    window.onload = function () {
        const sess = sessionStorage.getItem("user-admin");
        if (!sess) {
            window.location.href = "login.html";
            return;
        }
        
        // Initialize dashboard
        initializeDashboard();
    };

    // ===========================================
    // DASHBOARD INITIALIZATION
    // ===========================================
    function initializeDashboard() {
        setupEventListeners();
        loadOrders();
        showNotification('Dashboard admin dimuat', 'success');
    }

    // ===========================================
    // EVENT LISTENERS SETUP
    // ===========================================
    function setupEventListeners() {
        // Logout button
        document.getElementById("btn-logout").addEventListener("click", handleLogout);
        
        // Refresh button (manual only)
        const refreshBtn = document.querySelector(".btn-refresh");
        if (refreshBtn) {
            refreshBtn.addEventListener("click", () => {
                loadOrders();
                showNotification('Data pesanan diperbarui', 'info');
            });
        }
        
        // Tab switching
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                const tabName = btn.getAttribute('data-tab');
                switchTab(tabName);
            });
        });
        
        // Modal close
        const closeModal = document.querySelector('.close-modal');
        if (closeModal) {
            closeModal.addEventListener('click', closeModalHandler);
        }
        const overlay = document.getElementById('overlay');
        if (overlay) {
            overlay.addEventListener('click', closeModalHandler);
        }
    }

    // ===========================================
    // LOGOUT FUNCTIONALITY
    // ===========================================
    function handleLogout(e) {
        e.preventDefault();
        
        fetch(`${API_BASE}/admin/logout`, {
            method: "POST",
            credentials: "include"
        })
        .then(response => {
            sessionStorage.removeItem("user-admin");
            window.location.href = "login.html";
        })
        .catch(error => {
            console.log("Error logout:", error);
            // Force logout even if request fails
            sessionStorage.removeItem("user-admin");
            window.location.href = "login.html";
        });
    }

    // ===========================================
    // DATA LOADING FUNCTIONS
    // ===========================================
    function loadOrders() {
        showLoading(true);
        Promise.all([
            loadOrdersByStatus('pending'),
            loadOrdersByStatus('completed')
        ])
        .then(() => {
            showLoading(false);
            updateTabCounts();
        })
        .catch(error => {
            console.error('Error loading orders:', error);
            showLoading(false);
            showNotification('Gagal memuat data pesanan', 'error');
        });
    }

    // ===========================================
    // ORDERS LOADING BY STATUS
    // ===========================================
    function loadOrdersByStatus(status) {
        const statusParam = status === 'pending' ? 'process' : 'completed';
        
        return fetch(`${API_BASE}/admin/orders?status=${statusParam}`, {
            method: 'GET',
            credentials: 'include'
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to load ${status} orders`);
            }
            return response.json();
        })
        .then(data => {
            displayOrders(data.orders || [], status);
            return data.orders || [];
        })
        .catch(error => {
            console.error(`Error loading ${status} orders:`, error);
            showNotification(`Gagal memuat pesanan ${status}`, 'error');
            return [];
        });
    }

    // ===========================================
    // DISPLAY ORDERS
    // ===========================================
    function displayOrders(orders, status) {
        const containerId = status === 'pending' ? 'pendingOrders' : 'completedOrders';
        const container = document.getElementById(containerId);
        
        if (!container) {
            console.error(`Container ${containerId} not found`);
            return;
        }
        
        if (!orders || orders.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-inbox"></i>
                    <h3>Tidak ada pesanan ${status === 'pending' ? 'menunggu' : 'selesai'}</h3>
                    <p>${status === 'pending' ? 'Belum ada pesanan yang perlu diproses' : 'Belum ada pesanan yang selesai hari ini'}</p>
                </div>
            `;
            return;
        }

        container.innerHTML = orders.map(order => createOrderCard(order, status)).join('');
        
        // Add event listeners to action buttons
        container.querySelectorAll('.btn-action').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const orderId = e.target.getAttribute('data-order-id');
                const action = e.target.getAttribute('data-action');
                
                if (action === 'complete') {
                    completeOrder(orderId);
                } else if (action === 'detail') {
                    showOrderDetail(orderId, orders.find(o => o.id_order === orderId));
                }
            });
        });
    }

    // ===========================================
    // CREATE ORDER CARD HTML
    // ===========================================
    function createOrderCard(order, status) {
        const orderTime = new Date(order.created_at).toLocaleString('id-ID');
        const totalItems = order.items ? order.items.length : 0;
        
        return `
            <div class="order-card ${status}">
                <div class="order-header">
                    <div class="order-info">
                        <h4>Order #${order.id_order}</h4>
                        <p class="order-time">
                            <i class="fas fa-clock"></i>
                            ${orderTime}
                        </p>
                    </div>
                    <div class="order-status">
                        <span class="status-badge ${status}">
                            ${status === 'pending' ? 'Menunggu' : 'Selesai'}
                        </span>
                    </div>
                </div>
                
                <div class="order-items">
                    <h5>
                        <i class="fas fa-utensils"></i>
                        Items (${totalItems})
                    </h5>
                    <div class="items-list">
                        ${order.items ? order.items.slice(0, 3).map(item => 
                            `<span class="item-tag">${item.nama_produk} x${item.quantity}</span>`
                        ).join('') : ''}
                        ${totalItems > 3 ? `<span class="item-tag more">+${totalItems - 3} lainnya</span>` : ''}
                    </div>
                </div>
                
                <div class="order-total">
                    <strong>Total: Rp ${formatCurrency(order.total_harga)}</strong>
                </div>
                
                <div class="order-actions">
                    <button class="btn-action btn-detail" data-order-id="${order.id_order}" data-action="detail">
                        <i class="fas fa-eye"></i>
                        Detail
                    </button>
                    ${status === 'pending' ? `
                        <button class="btn-action btn-complete" data-order-id="${order.id_order}" data-action="complete">
                            <i class="fas fa-check"></i>
                            Tandai Selesai
                        </button>
                    ` : ''}
                </div>
            </div>
        `;
    }

    // ===========================================
    // COMPLETE ORDER
    // ===========================================
    function completeOrder(orderId) {
        if (!confirm('Apakah Anda yakin pesanan ini sudah selesai disiapkan?')) {
            return;
        }

        showLoading(true);
        
        fetch(`${API_BASE}/admin/orders/status`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
                order_id: orderId,
                status: 'completed'
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to complete order');
            }
            return response.json();
        })
        .then(data => {
            showLoading(false);
            showNotification('Pesanan berhasil ditandai selesai!', 'success');
            playSuccessSound();
            loadOrders(); // Refresh all data
        })
        .catch(error => {
            showLoading(false);
            console.error('Error completing order:', error);
            showNotification('Gagal menandai pesanan selesai', 'error');
        });
    }

    // ===========================================
    // SHOW ORDER DETAIL MODAL
    // ===========================================
    function showOrderDetail(orderId, orderData) {
        const modal = document.getElementById('orderModal');
        const overlay = document.getElementById('overlay');
        
        if (!modal || !overlay) {
            console.error('Modal elements not found');
            return;
        }
        
        // Show modal
        modal.classList.add('active');
        overlay.classList.add('active');
        
        // Populate modal with order details
        const modalContent = document.getElementById('orderDetailContent');
        if (modalContent && orderData) {
            modalContent.innerHTML = `
                <div class="order-detail">
                    <h4>Order #${orderData.id_order}</h4>
                    <p><strong>Waktu:</strong> ${new Date(orderData.created_at).toLocaleString('id-ID')}</p>
                    <p><strong>Status Customer:</strong> ${orderData.status_customer}</p>
                    <p><strong>Status Admin:</strong> ${orderData.status_admin}</p>
                    <p><strong>Total:</strong> Rp ${formatCurrency(orderData.total_harga)}</p>
                    
                    <h5>Item Pesanan:</h5>
                    <div class="order-items-detail">
                        ${orderData.items ? orderData.items.map(item => `
                            <div class="item-detail">
                                <span>${item.nama_produk}</span>
                                <span>x${item.quantity}</span>
                                <span>Rp ${formatCurrency(item.harga * item.quantity)}</span>
                            </div>
                        `).join('') : '<p>Tidak ada item</p>'}
                    </div>
                </div>
            `;
        }
    }

    // ===========================================
    // CLOSE MODAL
    // ===========================================
    function closeModalHandler() {
        const modal = document.getElementById('orderModal');
        const overlay = document.getElementById('overlay');
        
        if (modal) modal.classList.remove('active');
        if (overlay) overlay.classList.remove('active');
    }

    // ===========================================
    // TAB SWITCHING
    // ===========================================
    function switchTab(tabName) {
        currentTab = tabName;
        
        // Update tab buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
            if (btn.getAttribute('data-tab') === tabName) {
                btn.classList.add('active');
            }
        });
        
        // Update tab content
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
            if (content.id === `${tabName}Tab`) {
                content.classList.add('active');
            }
        });
    }

    // ===========================================
    // UPDATE TAB COUNTS
    // ===========================================
    function updateTabCounts() {
        const pendingContainer = document.getElementById('pendingOrders');
        const completedContainer = document.getElementById('completedOrders');
        
        if (pendingContainer) {
            const pendingCount = pendingContainer.querySelectorAll('.order-card').length;
            const pendingBadge = document.getElementById('pendingBadge');
            if (pendingBadge) pendingBadge.textContent = pendingCount;
        }
        
        if (completedContainer) {
            const completedCount = completedContainer.querySelectorAll('.order-card').length;
            const completedBadge = document.getElementById('completedBadge');
            if (completedBadge) completedBadge.textContent = completedCount;
        }
    }

    // ===========================================
    // UTILITY FUNCTIONS
    // ===========================================
    function showLoading(show) {
        const loading = document.getElementById('loading');
        if (loading) {
            if (show) {
                loading.classList.remove('hidden');
            } else {
                loading.classList.add('hidden');
            }
        }
    }

    function showNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.innerHTML = `
            <i class="fas fa-${type === 'success' ? 'check' : type === 'error' ? 'times' : 'info'}"></i>
            <span>${message}</span>
        `;
        
        document.body.appendChild(notification);
        
        // Auto remove after 3 seconds
        setTimeout(() => {
            notification.remove();
        }, 3000);
    }

    function formatCurrency(amount) {
        return new Intl.NumberFormat('id-ID').format(amount);
    }

    function playSuccessSound() {
        try {
            const sound = document.getElementById('successSound');
            if (sound) {
                sound.play();
            }
        } catch (error) {
            console.log('Could not play sound:', error);
        }
    }
}