import * as React from "react";
import {IconChevronDown, IconChevronUp} from "@tabler/icons-react";

export function SortableHeader({
                                   columnId,
                                   currentOrderBy,
                                   currentOrderDir,
                                   onSortChange,
                                   children,
                               }: {
    columnId: string
    currentOrderBy: string
    currentOrderDir: "asc" | "desc"
    onSortChange: (orderBy: string, orderDir: "asc" | "desc") => void
    children: React.ReactNode
}) {
    const isActive = currentOrderBy === columnId

    const nextOrderDir = isActive ? (currentOrderDir === "asc" ? "desc" : "asc") : "asc"

    return (
        <button
            type="button"
            onClick={() => onSortChange(columnId, nextOrderDir)}
            className="inline-flex items-center gap-1 font-semibold hover:underline"
        >
            {children}
            <SortIcon active={isActive} direction={currentOrderDir}/>
        </button>
    )
}

function SortIcon({active, direction}: { active: boolean; direction: "asc" | "desc" }) {
    if (!active) {
        return (
            <svg className="h-3 w-3 opacity-40" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                <path d="M6 9l6-6 6 6M6 15l6 6 6-6"/>
            </svg>
        )
    }
    if (direction === "asc") {
        return <IconChevronDown className="h-4 w-4"/>
    } else {
        return <IconChevronUp className="h-4 w-4"/>
    }
}
