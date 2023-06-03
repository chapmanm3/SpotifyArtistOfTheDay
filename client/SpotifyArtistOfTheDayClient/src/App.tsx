import { useState } from 'react'
import reactLogo from './assets/react.svg'
import './App.css'
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"

import getUserInfo from './api/getUserInfo'
import Header from './components/Header/header'

const queryClient = new QueryClient();

const AppWrappedWProviders = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  )
}

function App() {
  const [count, setCount] = useState(0)

  const handleLoginClick = () => {
    window.localStorage.setItem("is_authed", "true")
    window.location.href = "http://localhost:8080/api/login"
  }

  return (
    <div className="App">
      <Header />
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src="/vite.svg" className="logo" alt="Vite logo" />
        </a>
        <a href="https://reactjs.org" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <button
          onClick={handleLoginClick}
        >
          Test Login
        </button>
        <button
          onClick={() => getUserInfo()}>
          Get User Info
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  )
}

export default AppWrappedWProviders
