"use client"

import * as React from "react"
import axios from "axios"
import { QueryClient, useQuery, useQueryClient } from "@tanstack/react-query"
import { ColumnDef } from "@tanstack/react-table"
import { IconDotsVertical, IconPlus } from "@tabler/icons-react"
import { useIsMobile } from "@/hooks/use-mobile"
import { useLocalStorage } from "@/hooks/use-local-storage"
import { BACKEND_URL } from "@/constants"
import { Button } from "@/components/ui/button"
import {
    Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter
} from "@/components/ui/dialog"
import {
    DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import {
    AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
    AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { DataTable } from "@/components/ui/datatable"
import showErrorToast, { showSuccessToast } from "@/components/toast"
import { ScrollArea } from "@/components/ui/scroll-area"
import { SortableHeader } from "@/components/ui/sortable-header"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Link from "next/link"
import { Checkbox } from "@/components/ui/checkbox"
import {
    GroupModifierSyrve,
    Product,
    ProductVariation,
    ProductVariationGroup,
    SyrveProduct
} from "./types"
import { Category } from "@/app/dashboard/categories/types"
import { Pagination } from "@/app/dashboard/types"

async function fetchSyrveProducts(token?: string) {
    const res = await axios.get(`${BACKEND_URL}/syrve/products`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
    return res.data.data as SyrveProduct[]
}

const defaultFormData: Omit<Product, "id"> = {
    name: "",
    price: 0,
    sort: 0,
    category_id: "",
    image: null,
    external_id: "",
    sold: false,
    slug: "",
    hidden: false,
    alcohol: false,
    description: null,
    weight: 0,
}

export default function ProductsPage() {
    const [data, setData] = React.useState<Product[]>([])
    const [categories, setCategories] = React.useState<Category[]>([])
    const [selectedGroupModifiers, setSelectedGroupModifiers] = React.useState<(GroupModifierSyrve & {
        show: boolean
    })[]>([])
    const [pagination, setPagination] = React.useState<Pagination>({
        page: 1,
        total: 0,
        limit: 10,
        orderBy: "sort",
        orderDir: "asc",
    })
    const [editingItem, setEditingItem] = React.useState<Product | null>(null)
    const [dialogOpen, setDialogOpen] = React.useState(false)
    const [alertOpen, setAlertOpen] = React.useState(false)
    const [deleteId, setDeleteId] = React.useState<string | null>(null)
    const [formData, setFormData] = React.useState<Omit<Product, "id">>(defaultFormData)
    const [user] = useLocalStorage<any>("auth", null)
    const isMobile = useIsMobile()
    const queryClient = useQueryClient()

    React.useEffect(() => {
        axios.get(`${BACKEND_URL}/categories`)
            .then(res => setCategories(res.data.data))
            .catch(() => showErrorToast("Failed to load categories"))
    }, [])

    const { data: syrveData } = useQuery({
        queryKey: ["syrveProducts"],
        queryFn: () => fetchSyrveProducts(user?.access.token),
    })

    React.useEffect(() => {
        if (!editingItem) {
            setSelectedGroupModifiers([])
            return
        }

        async function loadVariationGroups() {
            try {
                const res = await axios.get(`${BACKEND_URL}/variation-groups?product_id=${editingItem.id}`)
                const groupModifiers = res.data.data as ProductVariationGroup[]
                const product = syrveData?.find(p => p.id === editingItem.external_id)
                const result = product?.groupModifiers?.map(gm => ({
                    ...gm,
                    show: groupModifiers.find(srv => srv.external_id === gm.id)?.show ?? false,
                })) ?? []
                setSelectedGroupModifiers(result)
            } catch {
                showErrorToast("Failed to load variation groups")
            }
        }

        loadVariationGroups()
    }, [editingItem, syrveData])

    const fetchData = React.useCallback(async () => {
        try {
            const res = await axios.get(`${BACKEND_URL}/products`, {
                params: {
                    page: pagination.page,
                    limit: pagination.limit,
                    order_by: pagination.orderBy,
                    order_dir: pagination.orderDir,
                },
            })
            setData(res.data.data)
            setPagination({
                page: res.data.pagination.page,
                total: res.data.pagination.total_count,
                limit: res.data.pagination.limit,
                orderBy: res.data.pagination.order_by,
                orderDir: res.data.pagination.order_dir,
            })
        } catch {
            showErrorToast("Failed to fetch products")
        }
    }, [pagination.page, pagination.limit, pagination.orderBy, pagination.orderDir])

    React.useEffect(() => {
        fetchData()
    }, [fetchData])

    React.useEffect(() => {
        if (editingItem) {
            const { id, ...rest } = editingItem
            setFormData(rest)
        }
    }, [editingItem])

    function openCreateDialog() {
        setEditingItem(null)
        setFormData(defaultFormData)
        setDialogOpen(true)
    }

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        const groupModifiersToShow = selectedGroupModifiers.filter(p => p.show)
        const regularModifiers = syrveData.find(p => p.id === formData.external_id).modifiers
        const requiredRegularModifiers = regularModifiers.filter(m => m.required)

        try {
            const method = editingItem ? "put" : "post"
            const url = editingItem
                ? `${BACKEND_URL}/products/${editingItem.id}`
                : `${BACKEND_URL}/products`

            const body = new FormData()
            Object.entries(formData).forEach(([key, value]) => {
                if (value !== null) {
                    if (key === "image" && value instanceof File) {
                        body.append("image", value)
                    } else {
                        body.append(key, String(value))
                    }
                }
            })

            const savedProductResponse = await axios[method](url, body, {
                headers: { Authorization: `Bearer ${user?.access.token}` },
            })

            const savedProduct: Product = savedProductResponse.data.data

            if (editingItem) {
                const delRes = await axios.delete(`${BACKEND_URL}/products/${editingItem.id}/variation-groups`, {
                    headers: { Authorization: `Bearer ${user?.access.token}` },
                })
                if (!delRes.data.success) throw new Error("Failed to delete variations")
            }

            for (const rm of requiredRegularModifiers) {
                const groupPayload: ProductVariationGroup = {
                    id: null,
                    product_id: savedProduct.id,
                    name: rm.name,
                    external_id: rm.id,
                    default_value: rm.defaultAmount,
                    show: false,
                    required: true,
                }
                const groupRes = await axios.post(`${BACKEND_URL}/variation-groups`, groupPayload, {
                    headers: { Authorization: `Bearer ${user?.access.token}` },
                })
            }

            for (const gm of groupModifiersToShow) {
                const groupPayload: ProductVariationGroup = {
                    id: null,
                    product_id: savedProduct.id,
                    name: gm.name,
                    external_id: gm.id,
                    default_value: null,
                    show: gm.show,
                    required: gm.required,
                }
                const groupRes = await axios.post(`${BACKEND_URL}/variation-groups`, groupPayload, {
                    headers: { Authorization: `Bearer ${user?.access.token}` },
                })
                const savedGroup: ProductVariationGroup = groupRes.data.data
                for (const variation of gm.childModifiers ?? []) {
                    const variationPayload: ProductVariation = {
                        id: null,
                        group_id: savedGroup.id!,
                        external_id: variation.id,
                        default_value: null,
                        show: true,
                        name: variation.name,
                    }
                    await axios.post(`${BACKEND_URL}/variations`, variationPayload, {
                        headers: { Authorization: `Bearer ${user?.access.token}` },
                    })
                }
            }

            fetchData()
            setDialogOpen(false)
            showSuccessToast(editingItem ? "Product updated" : "Product created")
        } catch (error: any) {
            const errorMessage = error.response?.data?.error || (editingItem ? "Failed to update product" : "Failed to create product")
            showErrorToast(errorMessage)
        }
    }

    async function handleDelete() {
        if (!deleteId) return
        try {
            await axios.delete(`${BACKEND_URL}/products/${deleteId}`, {
                headers: { Authorization: `Bearer ${user?.access.token}` },
            })
            fetchData()
            setDeleteId(null)
            setAlertOpen(false)
            showSuccessToast("Product deleted")
        } catch (error: any) {
            const errorMessage = error.response?.data?.error || "Failed to delete product"
            showErrorToast(errorMessage)
        }
    }

    function handleSortChange(orderBy: string, orderDir: "asc" | "desc") {
        setPagination((prev) => ({ ...prev, orderBy, orderDir, page: 1 }))
    }

    function handlePaginationChange(page: number, limit: number) {
        setPagination((prev) => ({ ...prev, page, limit }))
    }

    function resetSyrveCache() {
        queryClient.invalidateQueries({ queryKey: ["syrveProducts"] })
    }

    const columns: ColumnDef<Product>[] = [
        {
            accessorKey: "name",
            header: () => (
                <SortableHeader
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    columnId="name"
                    onSortChange={handleSortChange}
                >
                    Name
                </SortableHeader>
            ),
        },
        {
            accessorKey: "price",
            header: () => (
                <SortableHeader
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    columnId="price"
                    onSortChange={handleSortChange}
                >
                    Price
                </SortableHeader>
            ),
        },
        {
            accessorKey: "category_id",
            header: "Category",
            cell: ({ row }) => categories.find(c => c.id === row.original.category_id)?.name || "-",
        },
        {
            accessorKey: "slug",
            header: () => (
                <SortableHeader
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    columnId="slug"
                    onSortChange={handleSortChange}
                >
                    Slug
                </SortableHeader>
            ),
        },
        {
            accessorKey: "sort",
            header: () => (
                <SortableHeader
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    columnId="sort"
                    onSortChange={handleSortChange}
                >
                    Sort
                </SortableHeader>
            ),
        },
        {
            id: "actions",
            cell: ({ row }) => (
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="flex size-8">
                            <IconDotsVertical />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuItem
                            onClick={() => {
                                setEditingItem(row.original)
                                setDialogOpen(true)
                            }}
                        >
                            Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem
                            className="text-destructive"
                            onClick={() => {
                                setDeleteId(row.original.id)
                                setAlertOpen(true)
                            }}
                        >
                            Delete
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            ),
        },
    ]

    return (
        <>
            <DataTable
                data={data}
                columns={columns}
                pagination={pagination}
                customButtons={[
                    <Button
                        key="create"
                        variant="outline"
                        size="sm"
                        className="flex items-center text-white border-white hover:bg-white/10"
                        onClick={openCreateDialog}
                    >
                        <IconPlus />
                        <span className="hidden lg:inline">Create</span>
                        <span className="lg:hidden leading-none">Create</span>
                    </Button>,
                    <Button variant="outline" size="sm" onClick={resetSyrveCache}>
                        Reset Syrve Cache
                    </Button>,
                ]}
                onPaginationChange={handlePaginationChange}
            />

            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                <DialogContent
                    className={`h-full lg:h-auto w-full max-w-md shadow-lg ${isMobile ? "max-w-full" : ""}`}
                >
                    <DialogHeader>
                        <DialogTitle>{editingItem ? "Edit Product" : "New Product"}</DialogTitle>
                    </DialogHeader>

                    <ScrollArea className="h-[calc(100vh-10rem)]">
                        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                            <div className="flex flex-col gap-3">
                                <Label>Name</Label>
                                <Input
                                    value={formData.name ?? ""}
                                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                                />
                            </div>

                            <div className="flex flex-col gap-3">
                                <Label>Slug</Label>
                                <Input
                                    value={formData.slug ?? ""}
                                    onChange={(e) => setFormData(prev => ({ ...prev, slug: e.target.value }))}
                                />
                            </div>

                            <div className="flex flex-col gap-3">
                                <Label>Category</Label>
                                <select
                                    className="w-full border rounded p-2"
                                    value={formData.category_id}
                                    onChange={(e) => setFormData(prev => ({ ...prev, category_id: e.target.value }))}
                                >
                                    <option value="">Select Category</option>
                                    {categories.map(cat => (
                                        <option key={cat.id} value={cat.id}>{cat.name}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="flex flex-col gap-3">
                                <Label>
                                    Image {editingItem && typeof formData.image === "string" ? (
                                        <Link href={`${BACKEND_URL}/files/${formData.image}`} target="_blank">(view)</Link>
                                    ) : ""}
                                </Label>
                                <Input
                                    type="file"
                                    accept="image/*"
                                    onChange={e => {
                                        const file = e.target.files?.[0]
                                        if (file) setFormData(prev => ({ ...prev, image: file }))
                                    }}
                                />
                            </div>

                            {syrveData && (
                                <div className="flex flex-col gap-3">
                                    <Label>External ID</Label>
                                    <select
                                        className="w-full border rounded p-2"
                                        value={formData.external_id}
                                        onChange={e => {
                                            const val = e.target.value
                                            setFormData(prev => ({ ...prev, external_id: val }))
                                            const product = syrveData.find(p => p.id === val)
                                            if (product?.groupModifiers?.length) {
                                                setSelectedGroupModifiers(product.groupModifiers.map(gm => ({
                                                    ...gm,
                                                    show: false
                                                })))
                                            } else {
                                                setSelectedGroupModifiers([])
                                            }
                                        }}
                                    >
                                        <option value="">Select External ID</option>
                                        {syrveData.map(product => (
                                            <option key={product.id} value={product.id}>{product.name}</option>
                                        ))}
                                    </select>
                                </div>
                            )}

                            {selectedGroupModifiers.length > 0 && (
                                <div className="flex flex-col gap-2 mt-2">
                                    {selectedGroupModifiers.map((gm, idx) => (
                                        <div key={gm.id} className="flex items-center gap-2">
                                            <Input value={gm.name} disabled className="w-full" />
                                            <label className="flex items-center gap-1 select-none">
                                                <Checkbox
                                                    checked={gm.show}
                                                    onCheckedChange={() => {
                                                        setSelectedGroupModifiers(prev => {
                                                            const newMods = [...prev]
                                                            newMods[idx] = { ...newMods[idx], show: !newMods[idx].show }
                                                            return newMods
                                                        })
                                                    }}
                                                />
                                                Show
                                            </label>
                                        </div>
                                    ))}
                                </div>
                            )}

                            <div className="flex flex-col gap-3">
                                <Label>Description</Label>
                                <textarea
                                    className="border rounded px-3 py-2 h-24"
                                    value={formData.description ?? ""}
                                    onChange={e => setFormData(prev => ({ ...prev, description: e.target.value }))}
                                />
                            </div>

                            <div className="flex flex-col gap-3">
                                <Label>Price</Label>
                                <Input
                                    type="number"
                                    value={formData.price ?? 0}
                                    onChange={e => setFormData(prev => ({
                                        ...prev,
                                        price: parseFloat(e.target.value) || 0
                                    }))}
                                />
                            </div>

                            <div className="flex flex-col gap-3">
                                <Label>Weight</Label>
                                <Input
                                    type="number"
                                    value={formData.weight ?? 0}
                                    onChange={e => setFormData(prev => ({
                                        ...prev,
                                        weight: parseInt(e.target.value, 10) || 0
                                    }))}
                                />
                            </div>

                            <div className="flex flex-col gap-3">
                                <Label>Sort</Label>
                                <Input
                                    type="number"
                                    value={formData.sort ?? 0}
                                    onChange={(e) => setFormData((prev) => ({
                                        ...prev,
                                        sort: parseInt(e.target.value) || 0
                                    }))}
                                />
                            </div>


                            <div className="flex items-center gap-4">
                                <label className="flex items-center gap-1 select-none">
                                    <Checkbox
                                        checked={formData.sold}
                                        onCheckedChange={checked => setFormData(prev => ({
                                            ...prev,
                                            sold: Boolean(checked)
                                        }))}
                                    />
                                    Sold
                                </label>
                                <label className="flex items-center gap-1 select-none">
                                    <Checkbox
                                        checked={formData.hidden}
                                        onCheckedChange={checked => setFormData(prev => ({
                                            ...prev,
                                            hidden: Boolean(checked)
                                        }))}
                                    />
                                    Hidden
                                </label>
                                <label className="flex items-center gap-1 select-none">
                                    <Checkbox
                                        checked={formData.alcohol}
                                        onCheckedChange={checked => setFormData(prev => ({
                                            ...prev,
                                            alcohol: Boolean(checked)
                                        }))}
                                    />
                                    Alcohol
                                </label>
                            </div>

                            <DialogFooter className="bottom-0">
                                <Button className="w-full" type="submit">Submit</Button>
                            </DialogFooter>
                        </form>
                    </ScrollArea>
                </DialogContent>
            </Dialog>

            <AlertDialog open={alertOpen} onOpenChange={setAlertOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Are you sure you want to delete this product?</AlertDialogTitle>
                        <AlertDialogDescription>This action cannot be undone.</AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel onClick={() => setAlertOpen(false)}>Cancel</AlertDialogCancel>
                        <AlertDialogAction className="bg-destructive" onClick={handleDelete}>Delete</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}