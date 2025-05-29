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
            credentials: 'include', 
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

    window.onload = function () {
        const sess = sessionStorage.getItem("user");
        if (!sess) {
            window.location.href = "login.html";
        } else {
            showProduk();
        }
    }
}