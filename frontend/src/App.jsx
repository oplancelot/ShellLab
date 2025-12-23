import { useState } from 'react'
import './App.css'
import AtlasLootPage from './components/AtlasLootPage'
import DatabasePage from '../pages/DatabasePage/DatabasePage'
import SearchPage from './components/SearchPage'

function App() {
    const [activeTab, setActiveTab] = useState('atlas')

    return (
        <div className="app">
            <header className="header">
                <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
                    <h1>ShellLootLab</h1>
                    <nav style={{ display: 'flex', gap: '5px' }}>
                        <button 
                            onClick={() => setActiveTab('atlas')}
                            className={activeTab === 'atlas' ? 'active' : ''}
                            style={{ 
                                padding: '8px 16px', 
                                background: activeTab === 'atlas' ? '#404040' : '#2a2a2a',
                                border: 'none',
                                color: activeTab === 'atlas' ? '#fff' : '#888',
                                cursor: 'pointer',
                                borderRadius: '3px',
                                fontWeight: 'bold'
                            }}
                        >
                            AtlasLoot
                        </button>
                        <button 
                            onClick={() => setActiveTab('database')}
                            className={activeTab === 'database' ? 'active' : ''}
                            style={{ 
                                padding: '8px 16px', 
                                background: activeTab === 'database' ? '#404040' : '#2a2a2a',
                                border: 'none',
                                color: activeTab === 'database' ? '#fff' : '#888',
                                cursor: 'pointer',
                                borderRadius: '3px',
                                fontWeight: 'bold'
                            }}
                        >
                            Database
                        </button>
                        <button 
                            onClick={() => setActiveTab('search')}
                            className={activeTab === 'search' ? 'active' : ''}
                            style={{ 
                                padding: '8px 16px', 
                                background: activeTab === 'search' ? '#404040' : '#2a2a2a',
                                border: 'none',
                                color: activeTab === 'search' ? '#fff' : '#888',
                                cursor: 'pointer',
                                borderRadius: '3px',
                                fontWeight: 'bold'
                            }}
                        >
                            Search
                        </button>
                    </nav>
                </div>
            </header>

            <main style={{ flex: 1, overflow: 'hidden' }}>
                {activeTab === 'atlas' && <AtlasLootPage />}
                {activeTab === 'database' && <DatabasePage />}
                {activeTab === 'search' && <SearchPage />}
            </main>
        </div>
    )
}

export default App
