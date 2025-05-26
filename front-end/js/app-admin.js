// LOGIN ADMIN

const loginForm = document.getElementById("adminLoginForm");
if (loginForm) {
    loginForm.addEventListener("submit", function (e) {
        e.preventDefault();

        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;

        fetch("http://localhost:8080/api/admin/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
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
            })
    })
}

document.getElementById("btn-logout").addEventListener("click", function (e) {
    e.preventDefault();

    fetch("http://localhost:8080/api/admin/logout", {
        method: "POST",
        credentials: "include" // <--- INI WAJIB supaya cookie session ikut dikirim
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Logout gagal");
            }
            // Bersihkan sessionStorage
            sessionStorage.removeItem("user-admin");

            // Arahkan kembali ke halaman login
            window.location.href = "login.html";
        })
        .catch(error => {
            console.error("Gagal logout:", error);
            alert("Logout gagal: " + error.message);
        });
});