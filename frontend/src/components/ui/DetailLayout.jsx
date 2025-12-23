import React from 'react'
import { WowButton } from '../ui/Button'

/**
 * Shared layout for all detail views (Item, NPC, Quest, Spell, Faction)
 */
export const DetailPageLayout = ({ 
    children, 
    onBack,
    className = '' 
}) => (
    <div className={`flex-1 overflow-y-auto p-5 bg-bg-dark text-gray-200 ${className}`}>
        {onBack && (
            <WowButton 
                variant="back" 
                onClick={onBack}
                className="mb-6 group"
            >
                <span className="group-hover:-translate-x-1 transition-transform inline-block mr-2">←</span>
                Back to List
            </WowButton>
        )}
        {children}
    </div>
)

/**
 * Detail page header with icon and title
 */
export const DetailHeader = ({ 
    icon,
    iconBorderColor,
    title, 
    titleColor,
    subtitle,
    stats,
    children 
}) => (
    <header className="mb-8 pb-5 border-b border-border-dark">
        <div className="flex gap-5 items-start">
            {/* Icon */}
            {icon && (
                <div 
                    className="w-14 h-14 rounded border-2 shadow-lg overflow-hidden flex-shrink-0 bg-black/40 flex items-center justify-center"
                    style={{ borderColor: iconBorderColor || '#666' }}
                >
                    {icon}
                </div>
            )}
            
            {/* Title & Subtitle */}
            <div className="min-w-0 flex-1">
                <h1 
                    className="text-2xl font-bold m-0 leading-tight"
                    style={{ color: titleColor || '#fff' }}
                >
                    {title}
                </h1>
                {subtitle && (
                    <div className="text-gray-500 mt-1">{subtitle}</div>
                )}
                {stats && (
                    <div className="flex gap-4 mt-2 text-sm">
                        {stats}
                    </div>
                )}
            </div>
        </div>
        {children}
    </header>
)

/**
 * Section with gold header
 */
export const DetailSection = ({ title, children, className = '' }) => (
    <section className={`mb-8 ${className}`}>
        <h3 className="text-wow-gold font-bold uppercase text-sm tracking-wider mb-3 pb-1 border-b border-wow-gold/30">
            {title}
        </h3>
        {children}
    </section>
)

/**
 * Two-column grid for detail content
 */
export const DetailGrid = ({ children, className = '' }) => (
    <div className={`grid grid-cols-1 lg:grid-cols-2 gap-8 ${className}`}>
        {children}
    </div>
)

/**
 * Side panel (right column typically)
 */
export const DetailSidePanel = ({ children, className = '' }) => (
    <div className={`bg-bg-main p-5 rounded-lg border border-border-dark self-start ${className}`}>
        {children}
    </div>
)

/**
 * Loot/reward grid
 */
export const LootGrid = ({ children, className = '' }) => (
    <div className={`grid grid-cols-1 xl:grid-cols-2 gap-2 ${className}`}>
        {children}
    </div>
)

/**
 * Stat badge (e.g., "HP: 1000")
 */
export const StatBadge = ({ label, value, color }) => (
    <span 
        className="bg-black/30 px-2 py-0.5 rounded border border-white/5 text-sm"
        style={{ color: color || '#888' }}
    >
        {label}: <b className="text-gray-300">{value}</b>
    </span>
)

/**
 * Loading state for detail views
 */
export const DetailLoading = () => (
    <div className="flex-1 flex items-center justify-center bg-bg-dark">
        <div className="text-wow-gold italic animate-pulse text-lg">Loading...</div>
    </div>
)

/**
 * Error state for detail views
 */
export const DetailError = ({ message, onBack }) => (
    <div className="flex-1 flex flex-col items-center justify-center bg-bg-dark gap-4">
        <div className="text-red-500 font-bold text-lg">{message}</div>
        {onBack && (
            <WowButton variant="back" onClick={onBack}>← Back</WowButton>
        )}
    </div>
)

export default { 
    DetailPageLayout, 
    DetailHeader, 
    DetailSection, 
    DetailGrid, 
    DetailSidePanel, 
    LootGrid,
    StatBadge,
    DetailLoading,
    DetailError
}
