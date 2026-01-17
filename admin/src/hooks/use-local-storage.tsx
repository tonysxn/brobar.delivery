import {useEffect, useState} from "react";

export function useLocalStorage<T>(key: string, fallbackValue: T): [T, (value: T) => void] {
    const [value, setValue] = useState<T>(() => {
        if (typeof window === "undefined") return fallbackValue;
        const stored = localStorage.getItem(key);
        return stored ? JSON.parse(stored) : fallbackValue;
    });

    useEffect(() => {
        localStorage.setItem(key, JSON.stringify(value));
    }, [key, value]);

    return [value, setValue];
}
