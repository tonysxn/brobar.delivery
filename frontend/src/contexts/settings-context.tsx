"use client";

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';

export interface Setting {
    key: string;
    setting_type: string;
    value: string;
}

interface SettingsContextValue {
    settings: Setting[];
    loading: boolean;
    getSetting: (key: string) => Setting | null;
}

const SettingsContext = createContext<SettingsContextValue | null>(null);

const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL;

export function SettingsProvider({ children }: { children: ReactNode }) {
    const [settings, setSettings] = useState<Setting[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchSettings() {
            try {
                const res = await fetch(`${GATEWAY_URL}/settings`);
                if (res.ok) {
                    const json = await res.json();
                    setSettings(json.data || []);
                }
            } catch (error) {
                console.error("Error fetching settings:", error);
            } finally {
                setLoading(false);
            }
        }

        fetchSettings();
    }, []);

    const getSetting = (key: string): Setting | null => {
        return settings.find(s => s.key === key) || null;
    };

    return (
        <SettingsContext.Provider value={{ settings, loading, getSetting }}>
            {children}
        </SettingsContext.Provider>
    );
}

export function useSettings() {
    const context = useContext(SettingsContext);
    if (!context) {
        throw new Error('useSettings must be used within a SettingsProvider');
    }
    return context;
}
