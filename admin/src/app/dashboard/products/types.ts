export type Product = {
    id: string
    name: string
    price: number
    sort: number
    category_id: string
    image: File | string | null
    external_id: string
    sold: boolean
    slug: string
    hidden: boolean
    alcohol: boolean
    description: string | null
    weight: number
}

export interface Modifier {
    id: string
    name: string
    defaultAmount: number
    required: boolean
}

export interface GroupModifierSyrve {
    id: string
    name: string
    required: boolean
    defaultAmount: number
    childModifiers: Modifier[]
}

export interface ProductVariationGroup {
    id: string | null
    name: string
    required: boolean
    default_value: number | null
    product_id: string
    show: boolean
    external_id: string
}

export type ProductVariation = {
    id: string | null
    group_id: string
    external_id: string
    default_value: number | null
    show: boolean
    name: string
}

export interface SyrveProduct {
    id: string
    name: string
    modifiers: Modifier[]
    groupModifiers: GroupModifierSyrve[]
}