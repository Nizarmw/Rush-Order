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
    }

    function addMenuBehavior() {
        const AddMenuBtn = document.querySelectorAll('.add-to-cart-btn');
        AddMenuBtn.forEach(btn => {
            btn.addEventListener('click', () => {
                console.log('aw dipencet');
                animateCartIcon();
                showAlert("Berhasil Tambah Menu", "success")
            })

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
    

    window.onload = function () {
        const sess = sessionStorage.getItem("user");
        if (!sess) {
            window.location.href = "login.html";
        } else {
            showProduk();

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
}