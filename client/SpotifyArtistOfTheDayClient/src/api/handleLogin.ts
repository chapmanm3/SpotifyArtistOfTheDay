export default function handleLogin() {
  window.localStorage.setItem("is_authed", "true")
  window.location.href = "http://localhost:8080/api/login"
}

