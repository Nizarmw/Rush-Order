// LOGIN
document.querySelector("form").addEventListener('submit', function(e) {
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