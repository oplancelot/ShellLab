import { useState } from 'react'
import AtlasLootPage from './pages/AtlasLootPage/AtlasLootPage'
import DatabasePage from './pages/DatabasePage/DatabasePage'
import SearchPage from './pages/SearchPage/SearchPage'
import SettingsPage from './pages/SettingsPage'
import { TabButton } from './components/ui'

function App() {
    const [activeTab, setActiveTab] = useState('atlas')

    return (
        <div className="h-screen flex flex-col bg-bg-dark text-white">
            {/* Header */}
            <header className="bg-gradient-to-b from-[#2a2a3a] to-bg-main border-b-[3px] border-bg-dark px-5 py-3 flex items-center justify-between">
                <div className="flex items-center gap-5">
                    <h1 className="text-2xl text-wow-gold font-normal drop-shadow-md flex items-center gap-2.5">
                        <img src="/shelllab-logo.svg" alt="ShellLab" className="w-8 h-8" />
                        ShellLab
                    </h1>
                    <nav className="flex gap-1">
                        <TabButton 
                            active={activeTab === 'atlas'} 
                            onClick={() => setActiveTab('atlas')}
                        >
                            AtlasLoot
                        </TabButton>
                        <TabButton 
                            active={activeTab === 'database'} 
                            onClick={() => setActiveTab('database')}
                        >
                            Database
                        </TabButton>
                        <TabButton 
                            active={activeTab === 'search'} 
                            onClick={() => setActiveTab('search')}
                        >
                            Search
                        </TabButton>
                        <TabButton 
                            active={activeTab === 'settings'} 
                            onClick={() => setActiveTab('settings')}
                        >
                            Settings
                        </TabButton>
                    </nav>
                </div>
            </header>

            {/* Main Content */}
            <main className="flex-1 overflow-hidden">
                {activeTab === 'atlas' && <AtlasLootPage />}
                {activeTab === 'database' && <DatabasePage />}
                {activeTab === 'search' && <SearchPage />}
                {activeTab === 'settings' && <SettingsPage />}
            </main>
        </div>
    )
}

export default App
