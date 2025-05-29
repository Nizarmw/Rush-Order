// LOGIN

if (window.location.pathname.endsWith('login.html')) {
    document.querySelector("form").addEventListener('submit', function (e) {
        e.preventDefault(); // mencegah reload form

        // Ambil input nama dan meja
        const nama = document.getElementById('nama').value;
        const meja = document.getElementById('meja').value;

        // Kirim ke backend API
        fetch('http://localhost:8080/api/sessions/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include', // WAJIB
            body: JSON.stringify({
                nama: nama,
                meja: parseInt(meja)
            })
        })        

        .then(response => {
            if (!response.ok) {
                throw new Error("Gagal login");
            }
            return response.json();
        })

        .then(data => {
            console.log("Login berhasil:", data);
            // Contoh: simpan data ke localStorage/sessionStorage
            sessionStorage.setItem('user', JSON.stringify(data));

            // Redirect ke halaman menu (misalnya)
            window.location.href = 'index.html';
        })
        .catch(error => {
            console.log(error);
            alert('Login gagal: ' + error.message);
        });
    });
}

if (window.location.pathname.endsWith('index.html')) {

    function showProduk() {
        fetch('http://localhost:8080/api/produk/')
            .then(response => {
                if (!response.ok) {
                    throw new Error("Gagal Fetch Produk");
                }
                return response.json();
            })
            .then(data => {
                console.log("Hasil Fetch:", data);
                renderProduk(data.produk);
            })
            .catch(error => {
                console.error('Error ambil data produk', error);
            })
    }

    function renderProduk(listProduk) {
        const gridMenu = document.getElementById('menuGrid');

        gridMenu.innerHTML = listProduk.map(item => `
        <div class="menu-item fade-in">
            <div class="menu-item-image" style="background-image: url('${item.image_url}');">
            </div>
            <div class="menu-item-content">
                <h3 class="menu-item-title">${item.nama_produk}</h3>
                <p class="menu-item-description">${item.deskripsi}</p>
                <div class="menu-item-price">Rp ${item.harga_produk}</div>
                <div class="quantity-controls">
                    <button class="quantity-btn" onclick="changeQuantity('${item.id_produk}', -1)">-</button>
                    <input type="number" class="quantity-input" id="qty-${item.id_produk}" value="1" min="1" max="10">
                    <button class="quantity-btn" onclick="changeQuantity('${item.id_produk}', 1)">+</button>
                </div>
                <button class="add-to-cart-btn">
                    <i class="fas fa-plus"></i>
                    Tambah ke Keranjang
                </button>
            </div>
        </div>
    `).join('');
        addMenuBehavior();
    }    function addMenuBehavior() {
        const AddMenuBtn = document.querySelectorAll('.add-to-cart-btn');
        AddMenuBtn.forEach(btn => {
            btn.addEventListener('click', async () => {
                const inp_data = btn.parentElement.querySelector('.quantity-controls input')
                const idProduk = inp_data.id.substring(4, 10)
                const namaProduk = btn.parentElement.querySelector('h3').innerHTML
                const jumlah = parseInt(inp_data.value)
                const harga = parseInt(btn.parentElement.querySelector('.menu-item-price').innerHTML.replace(/[^0-9]/g, ""))
            
                try {
                    // Add to cart
                    const cartResponse = await fetch('http://localhost:8080/api/carts/', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        credentials: 'include',
                        body: JSON.stringify({
                            id_produk: idProduk,
                            nama_produk: namaProduk,
                            jumlah: jumlah,
                            harga: harga
                        })
                    });
            
                    const cartData = await cartResponse.json();
            
                    if (cartData.message) {
                        // Wait a bit then fetch session
                        await new Promise(resolve => setTimeout(resolve, 100));
                        
                        try {
                            const sessionResponse = await fetch('http://localhost:8080/api/sessions/', {
                                credentials: 'include'
                            });
            
                            console.log('Session response status:', sessionResponse.status);
                            
                            if (sessionResponse.ok) {
                                const userData = await sessionResponse.json();
                                console.log('Session data:', userData);
                                
                                sessionStorage.removeItem('user');
                                sessionStorage.setItem('user', JSON.stringify(userData));
                                animateCartIcon();
                                showAlert("Berhasil Tambah Menu", "success");
                            } else {
                                throw new Error(`Session fetch failed: ${sessionResponse.status}`);
                            }

                            loadCart();
                        } catch (sessionError) {
                            console.error("Session update failed:", sessionError);
                            // Cart berhasil, tapi session update gagal
                            animateCartIcon();
                            showAlert("Item ditambahkan, tapi gagal update session", "warning");
                        }
                    } else {
                        showAlert(cartData.error || "Gagal menambahkan", "error");
                    }
                } catch (err) {
                    console.error('Error:', err);
                    showAlert("Terjadi kesalahan saat menambahkan", "error");
                }
            });
        })
    }

    function animateCartIcon() {
        const cartIcon = document.querySelector('.cart-icon');
        cartIcon.classList.add('cart-pulse');
        setTimeout(() => {
            cartIcon.classList.remove('cart-pulse');
        }, 500);
    }

    function showAlert(message, type = 'info') {
        const alertClass = type === 'error' ? 'alert-error' : 'alert-success';

        const existingAlert = document.querySelector('.app-alert');
        if (existingAlert) {
            existingAlert.remove();
        }

        const alert = document.createElement('div');
        alert.className = `app-alert ${alertClass}`;
        alert.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: ${type === 'error' ? '#fee' : '#efe'};
            color: ${type === 'error' ? '#c53030' : '#2f855a'};
            padding: 15px 20px;
            border-radius: 8px;
            border-left: 4px solid ${type === 'error' ? '#e53e3e' : '#48bb78'};
            box-shadow: 0 4px 20px rgba(0,0,0,0.1);
            z-index: 10000;
            font-weight: 600;
            max-width: 300px;
            animation: slideInFromRight 0.3s ease-out;
        `;
        alert.textContent = message;

        document.body.appendChild(alert);

        setTimeout(() => {
            if (alert.parentNode) {
                alert.style.animation = 'slideOutToRight 0.3s ease-in';
                setTimeout(() => alert.remove(), 300);
            }
        }, 3000);
    }

    function changeQuantity(productId, change) {
        const input = document.getElementById(`qty-${productId}`);
        const currentValue = parseInt(input.value);
        const newValue = Math.max(1, Math.min(10, currentValue + change));
        input.value = newValue;
    }

    function logout() {
        fetch("http://127.0.0.1:8080/api/sessions/clear", {
            method: "POST",
            credentials: "include"
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Logout gagal");
            }
            return response.json();
        })
        .then(data => {
            sessionStorage.clear()
            console.log(data.message);
            window.location.href = "index.html";
        })
        .catch(error => {
            console.error("Error:", error);
            alert("Terjadi kesalahan saat logout.");
        });
    }

    function toggleCart() {
        const cartSidebar = document.getElementById('cartSidebar');
        const overlay = document.getElementById('overlay');
        
        const isOpen = cartSidebar.classList.contains('open');
        
        if (isOpen) {
            cartSidebar.classList.remove('open');
            overlay.classList.remove('active');
        } else {
            cartSidebar.classList.add('open');
            overlay.classList.add('active');
        }
    }

    async function paymentMidtrans() {
        try {
            const response = await fetch('http://localhost:8080/api/payment/checkout', {
                method: 'POST',
                credentials: 'include', // penting untuk mengirim session cookie
                headers: {
                    'Content-Type': 'application/json'
                }
            });
    
            const data = await response.json();
    
            if (!response.ok) {
                alert(data.error || "Gagal membuat pembayaran.");
                return;
            }
    
            const snapToken = data.snap_token;
    
            if (!snapToken) {
                alert("Token pembayaran tidak ditemukan.");
                return;
            }            // Panggil Snap popup Midtrans
            window.snap.pay(snapToken, {                onSuccess: function(result) {
                    console.log("Payment Success:", result);
                    
                    // Store order ID for status tracking
                    sessionStorage.setItem('currentOrderId', data.order_id);
                    
                    // Set flag that payment just completed successfully
                    sessionStorage.setItem('paymentJustCompleted', 'true');
                    
                    // Show success page
                    showSuccessPage(data.order_id, result);
                },
                onPending: function(result) {
                    console.log("Payment Pending:", result);
                    
                    // Store order ID for status tracking
                    sessionStorage.setItem('currentOrderId', data.order_id);
                    
                    // Show alert and redirect to status page
                    alert("Pembayaran masih pending. Silakan cek status pesanan Anda.");
                    checkOrderStatus();
                },
                onError: function(result) {
                    alert("Terjadi kesalahan saat pembayaran.");
                    console.error("Error:", result);
                },
                onClose: function() {
                    alert("Kamu menutup popup sebelum selesai.");
                }
            });
    
        } catch (err) {
            console.error("Terjadi kesalahan:", err);
            alert("Terjadi kesalahan saat memproses pembayaran.");
        }
    }   
    

    window.onload = function () {
        const sess = sessionStorage.getItem("user");
        if (!sess) {
            window.location.href = "login.html";
        } else {
            showProduk();

            document.getElementById('checkoutBtn').addEventListener('click', paymentMidtrans)

            fetch('http://localhost:8080/api/sessions/', {
                credentials: 'include'
            })
            .then(res => res.json())
            .then(userData => {
                console.log(userData)
            })
            .catch(err => {
                console.error("Gagal fetch session user:", err);
                showAlert("Item ditambahkan, tapi gagal update session", "warning");
            });
        }
    }

    window.delCart = async function(id) {
        try {
            const response = await fetch(`http://localhost:8080/api/carts/?id=${id}`, {
                method: 'DELETE',
                credentials: 'include'
            });
    
            const result = await response.json();
    
            if (response.ok) {
                showAlert(result.message || "Item dihapus dari cart", "success");
                animateCartIcon();
                try {
                    const sessionResponse = await fetch('http://localhost:8080/api/sessions/', {
                        credentials: 'include'
                    });
    
                    console.log('Session response status:', sessionResponse.status);
                    
                    if (sessionResponse.ok) {
                        const userData = await sessionResponse.json();
                        console.log('Session data:', userData);
                        
                        sessionStorage.removeItem('user');
                        sessionStorage.setItem('user', JSON.stringify(userData));
                        animateCartIcon();
                        showAlert("Berhasil Hapus Menu", "success");
                    } else {
                        throw new Error(`Session fetch failed: ${sessionResponse.status}`);
                    }

                    loadCart();
                } catch (sessionError) {
                    console.error("Session update failed:", sessionError);
                    // Cart berhasil, tapi session update gagal
                    animateCartIcon();
                    showAlert("Item ditambahkan, tapi gagal update session", "warning");
                }
            } else {
                showAlert(result.error || "Gagal menghapus item", "error");
            }
        } catch (err) {
            console.error("Delete error:", err);
            showAlert("Terjadi kesalahan saat menghapus", "error");
        }
    }    

    async function loadCart() {
        try {
            const response = await fetch('http://localhost:8080/api/carts/', {
                credentials: 'include'
            });
            if (!response.ok) throw new Error('Gagal mengambil cart');
            const data = await response.json();

            // Render items
            const cartItemsDiv = document.getElementById('cartItems');
            if (!cartItemsDiv) return;
            if (!data.items || Object.keys(data.items).length === 0) {
                cartItemsDiv.innerHTML = '<p>Keranjang kosong.</p>';
            } else {
                cartItemsDiv.innerHTML = Object.values(data.items).map(item => `
                    <div class="cart-item">
                        <div class="cart-item-info">
                            <h4>${item.nama_produk}</h4>
                            <p>Jumlah: ${item.jumlah}</p>
                            <p>Subtotal: Rp ${item.subtotal}</p>
                        </div>
                        <div class="cart-item-controls">
                        <button id="min-${item.id_produk}" class="quantity-btn" style="color: #e53e3e; margin-left: 10px;" onclick="delCart('${item.id_produk}')">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                    </div>
                `).join('');
            }

            // Render total
            const cartTotalSpan = document.getElementById('cartTotal');
            if (cartTotalSpan) cartTotalSpan.textContent = data.total || 0;

        } catch (err) {
            console.error('Error loading cart:', err);
        }
    }

    // Order Status Tracking Functions
    function showSuccessPage(orderId, paymentResult) {
        // Generate receipt content
        generateReceiptContent(orderId);
        
        // Show success page
        document.getElementById('menuPage').classList.remove('active');
        document.getElementById('checkoutPage').classList.remove('active');
        document.getElementById('orderStatusPage').classList.remove('active');
        document.getElementById('successPage').classList.add('active');
    }

    function generateReceiptContent(orderId) {
        const receiptDiv = document.getElementById('receiptContent');
        const user = JSON.parse(sessionStorage.getItem('user') || '{}');
        const now = new Date().toLocaleString('id-ID');
        
        receiptDiv.innerHTML = `
            Order ID: ${orderId}
            Nama: ${user.nama || '-'}
            Meja: ${user.meja || '-'}
            Waktu: ${now}
            
            Status: Pembayaran Berhasil
            Pesanan sedang diproses...
        `;
    }    window.checkOrderStatus = function() {
        const orderId = sessionStorage.getItem('currentOrderId');
        if (!orderId) {
            alert('Order ID tidak ditemukan. Silakan pesan kembali.');
            return;
        }
        
        showOrderStatusPage(orderId);
    }

    function showOrderStatusPage(orderId) {
        // Hide other pages
        document.getElementById('menuPage').classList.remove('active');
        document.getElementById('checkoutPage').classList.remove('active');
        document.getElementById('successPage').classList.remove('active');
        document.getElementById('orderStatusPage').classList.add('active');
        
        // Load order status
        loadOrderStatus(orderId);
        
        // Start auto-refresh for status updates (every 30 seconds)
        startAutoRefresh(orderId);
    }    async function loadOrderStatus(orderId) {
        try {
            showLoading(true);
            
            const response = await fetch(`http://localhost:8080/api/order/${orderId}/status`, {
                credentials: 'include'
            });
            
            if (!response.ok) {
                throw new Error('Gagal mengambil status pesanan');
            }
            
            const orderData = await response.json();
            
            // Clear payment success flag after first load
            setTimeout(() => {
                sessionStorage.removeItem('paymentJustCompleted');
            }, 3000);
            
            displayOrderStatus(orderData);
            
        } catch (error) {
            console.error('Error loading order status:', error);
            alert('Gagal memuat status pesanan: ' + error.message);
        } finally {
            showLoading(false);
        }
    }

    function displayOrderStatus(orderData) {
        // Update order info
        document.getElementById('statusOrderId').textContent = orderData.id_order;
        document.getElementById('statusOrderTotal').textContent = orderData.total_harga;
        
        // Update timeline based on status
        updateStatusTimeline(orderData.status_customer, orderData.status_admin);
        
        // Display order items
        displayOrderItems(orderData.items);
    }    function updateStatusTimeline(customerStatus, adminStatus) {
        // Reset all steps (hanya 3 step sekarang)
        const steps = ['step-success', 'step-process', 'step-completed'];
        steps.forEach(stepId => {
            const element = document.getElementById(stepId);
            element.classList.remove('active', 'completed');
        });
        
        // Check if payment just completed from Midtrans
        const paymentJustCompleted = sessionStorage.getItem('paymentJustCompleted');
        
        // Always start from success step if payment was made
        if (customerStatus === 'success' || paymentJustCompleted === 'true') {
            document.getElementById('step-success').classList.add('completed');
            
            // Update based on admin status
            if (adminStatus === 'process') {
                document.getElementById('step-process').classList.add('active');
            } else if (adminStatus === 'completed') {
                document.getElementById('step-process').classList.add('completed');
                document.getElementById('step-completed').classList.add('active');
            } else {
                // Payment success, waiting for admin to process
                document.getElementById('step-process').classList.add('active');
            }
        } else if (customerStatus === 'pending') {
            // If still pending but came from Midtrans success, treat as success
            if (paymentJustCompleted === 'true') {
                document.getElementById('step-success').classList.add('completed');
                document.getElementById('step-process').classList.add('active');
            } else {
                document.getElementById('step-success').classList.add('active');
            }
        } else {
            // Default to first step active
            document.getElementById('step-success').classList.add('active');
        }
    }

    function displayOrderItems(items) {
        const itemsContainer = document.getElementById('statusOrderItems');
        
        if (!items || items.length === 0) {
            itemsContainer.innerHTML = '<p>Tidak ada item dalam pesanan.</p>';
            return;
        }
        
        itemsContainer.innerHTML = items.map(item => `
            <div class="status-order-item">
                <div class="status-item-info">
                    <div class="status-item-name">${item.nama_produk || 'Item'}</div>
                    <div class="status-item-quantity">Jumlah: ${item.jumlah}</div>
                </div>
                <div class="status-item-price">Rp ${item.subtotal}</div>
            </div>
        `).join('');
    }

    window.refreshOrderStatus = function() {
        const orderId = sessionStorage.getItem('currentOrderId');
        if (orderId) {
            loadOrderStatus(orderId);
        }
    }

    // Auto-refresh functionality
    let autoRefreshInterval = null;
      function startAutoRefresh(orderId) {
        // Clear any existing interval
        stopAutoRefresh();
        
        // Start new interval (refresh every 30 seconds)
        autoRefreshInterval = setInterval(() => {
            loadOrderStatus(orderId);
        }, 30000);
    }
    
    function stopAutoRefresh() {
        if (autoRefreshInterval) {
            clearInterval(autoRefreshInterval);
            autoRefreshInterval = null;
        }
    }
    
    // Stop auto-refresh when leaving order status page
    window.goBackFromStatus = function() {
        stopAutoRefresh();
        
        // Check if we came from success page or have order ID
        const currentOrderId = sessionStorage.getItem('currentOrderId');
        
        if (currentOrderId) {
            // If we have an order ID, go back to success page instead of menu
            showSuccessPage(currentOrderId, null);
        } else {
            // No order ID, go to menu if user exists
            const user = sessionStorage.getItem('user');
            if (user) {
                showMenuPage();
            } else {
                window.location.href = 'login.html';
            }        }
    }

    window.goBackToMenu = function() {
        stopAutoRefresh();
        sessionStorage.removeItem('currentOrderId');
        sessionStorage.removeItem('paymentJustCompleted');
        showMenuPage();
    }

    window.orderAgain = function() {
        stopAutoRefresh();
        sessionStorage.removeItem('currentOrderId');
        sessionStorage.removeItem('paymentJustCompleted');
        showMenuPage();
    }

    function showMenuPage() {
        document.getElementById('successPage').classList.remove('active');
        document.getElementById('orderStatusPage').classList.remove('active');
        document.getElementById('checkoutPage').classList.remove('active');
        document.getElementById('menuPage').classList.add('active');
    }    function showLoading(show) {
        const loadingElement = document.getElementById('loading');
        if (show) {
            loadingElement.classList.remove('hidden');
        } else {
            loadingElement.classList.add('hidden');
        }
    }
}