import axios from 'axios'

axios.interceptors.response.use((resp) => {
  if (resp.status == 403) {
    window.localStorage.setItem("is_authed", "false")
  }
  return resp
})
