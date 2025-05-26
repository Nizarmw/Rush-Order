// LOGIN

if (window.location.pathname.endsWith('login.html')) {
    const loginForm = document.getElementById("adminLoginForm");
    loginForm.addEventListener("submit", function(e) {
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

                window.location.href="index.html";
            })
        .catch(error => {
            console.log(error);
            alert('Login gagal: ' + error.message);
        })
    })
    }

if (window.location.pathname.endsWith('index.html')) {
    document.getElementById("btn-logout").addEventListener("click", function (e) {
        e.preventDefault();

        // (Opsional) Panggil endpoint logout backend
        fetch("http://localhost:8080/api/admin/logout", {
            method: "POST",
            // Tidak perlu credentials karena tidak pakai cookie
            headers: {
                "Content-Type": "application/json"
            }
        })

        .then(response => {
            if (response.ok) {
                console.log("Logout berhasil");
                // Redirect ke login
                window.location.href = "login.html";
            } else {
                return response.json().then(data => {
                    console.error("Logout gagal:", data.error);
                });
            }
        })
        .catch(error => {
            console.log("Error logout:", error);
        });
    });

    window.onload = function () {
        const sess = sessionStorage.getItem("user-admin");
        if (!sess) {
            window.location.href = "login.html";
        }
    }
}