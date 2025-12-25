import React, { useState } from 'react'
import { FixMissingIcons } from '../../../wailsjs/go/main/App'
import { PageLayout } from '../../components/ui'

function SettingsPage() {
    const [fixing, setFixing] = useState(false)
    const [result, setResult] = useState(null)
    const [iconType, setIconType] = useState('item') // 'item' or 'spell'

    const handleFixIcons = async () => {
        setFixing(true)
        setResult(null)
        
        try {
            const res = await FixMissingIcons(iconType, 0)
            setResult(res)
        } catch (error) {
            setResult({
                totalMissing: 0,
                fixed: 0,
                failed: 0,
                message: `Error: ${error.toString()}`
            })
        } finally {
            setFixing(false)
        }
    }

    return (
        <PageLayout>
            <div className="p-8 max-w-4xl mx-auto">
                <h1 className="text-3xl font-bold text-white mb-8">Settings & Tools</h1>

                {/* Icon Fix Section */}
                <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 mb-6">
                    <h2 className="text-xl font-bold text-wow-gold mb-4">Fix Missing Icons</h2>
                    
                    <p className="text-gray-300 mb-4">
                        Automatically fetch icon names from the Turtle WoW database and update the local database.
                    </p>

                    {/* Icon Type Selector */}
                    <div className="mb-4">
                        <label className="text-gray-300 text-sm mb-2 block">Icon Type:</label>
                        <div className="flex gap-4">
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="radio"
                                    name="iconType"
                                    value="item"
                                    checked={iconType === 'item'}
                                    onChange={(e) => setIconType(e.target.value)}
                                    className="w-4 h-4"
                                />
                                <span className="text-white">Item Icons</span>
                            </label>
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="radio"
                                    name="iconType"
                                    value="spell"
                                    checked={iconType === 'spell'}
                                    onChange={(e) => setIconType(e.target.value)}
                                    className="w-4 h-4"
                                />
                                <span className="text-white">Spell Icons</span>
                            </label>
                        </div>
                    </div>

                    <button
                        onClick={handleFixIcons}
                        disabled={fixing}
                        className={`
                            px-6 py-3 rounded font-bold text-sm uppercase tracking-wider
                            transition-all duration-200
                            ${fixing 
                                ? 'bg-gray-600 text-gray-400 cursor-not-allowed' 
                                : 'bg-wow-gold hover:bg-yellow-500 text-gray-900 hover:shadow-lg'
                            }
                        `}
                    >
                        {fixing ? (
                            <span className="flex items-center gap-2">
                                <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none"/>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
                                </svg>
                                Fixing Icons...
                            </span>
                        ) : (
                            `Fix All Missing ${iconType === 'item' ? 'Item' : 'Spell'} Icons`
                        )}
                    </button>

                    {/* Result Display */}
                    {result && (
                        <div className={`mt-4 p-4 rounded border ${
                            result.fixed > 0 
                                ? 'bg-green-900/20 border-green-500/50 text-green-200' 
                                : result.totalMissing === 0
                                    ? 'bg-blue-900/20 border-blue-500/50 text-blue-200'
                                    : 'bg-red-900/20 border-red-500/50 text-red-200'
                        }`}>
                            <div className="font-bold mb-2">{result.message}</div>
                            <div className="text-sm space-y-1 opacity-90">
                                <div>Total Missing: {result.totalMissing}</div>
                                <div>Fixed: {result.fixed}</div>
                                {result.failed > 0 && <div>Failed: {result.failed}</div>}
                            </div>
                            {result.totalMissing > result.fixed && result.fixed > 0 && (
                                <div className="mt-3 text-sm">
                                    Click the button again to fix more icons.
                                </div>
                            )}
                        </div>
                    )}

                    <div className="mt-4 text-sm text-gray-500">
                        <div className="font-semibold mb-1">How it works:</div>
                        <ul className="list-disc list-inside space-y-1 text-xs">
                            <li>Fetches icon data from https://database.turtlecraft.gg</li>
                            <li>Updates database with correct icon names</li>
                            <li>Icon images are auto-downloaded by the icon service</li>
                            <li>Process 100 items at a time to avoid server load</li>
                        </ul>
                    </div>
                </div>

                {/* About Section */}
                <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6">
                    <h2 className="text-xl font-bold text-wow-gold mb-4">About ShellLab</h2>
                    <div className="text-gray-300 text-sm space-y-2">
                        <p>ShellLab - WoW Toolkit for Turtle WoW</p>
                        <p className="text-gray-500">Browse items, quests, NPCs, and more from the Turtle WoW database.</p>
                    </div>
                </div>
            </div>
        </PageLayout>
    )
}

export default SettingsPage
