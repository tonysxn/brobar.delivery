
export function getAssetUrl(path: string | undefined | null): string {
    if (!path) return "";

    // If it's already a full URL, return it
    if (path.startsWith('http')) return path;

    const baseUrl = process.env.NEXT_PUBLIC_FILE_URL || "http://localhost:3001";

    // Handle cases where path might already have a leading slash
    const cleanPath = path.startsWith('/') ? path.slice(1) : path;

    return `${baseUrl}/${cleanPath}`;
}
