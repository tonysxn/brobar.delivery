export type Pagination = {
    page: number
    total: number
    limit: number
    orderBy: string
    orderDir: "asc" | "desc"
}