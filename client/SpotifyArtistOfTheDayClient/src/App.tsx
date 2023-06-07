import './App.css'
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"

import Header from './components/Header/header'
import initInterceptors from './api/axios/initInterceptors'

const queryClient = new QueryClient();

const AppWrappedWProviders = () => {
  initInterceptors()
  return (
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  )
}

function App() {
  return (
    <div className="App">
      <div className='header-container'>
        <Header />
      </div>
    </div>
  )
}

export default AppWrappedWProviders
