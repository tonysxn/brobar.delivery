"use client"

import { GoogleMap, useJsApiLoader, Polygon, Marker } from "@react-google-maps/api";
import { useState, useCallback, useMemo, useEffect, useRef } from "react";
import { toast } from "sonner";
import { useIsMobile } from "@/hooks/use-mobile";
import { useSettings } from "@/contexts/settings-context";
import { useDeliveryZones } from "@/hooks/use-delivery-zones";

const containerStyle = {
    width: '100%',
    height: '500px',
    borderRadius: '1rem'
};

const defaultCenter = {
    lat: 50.0014656,
    lng: 36.245192
};

const libraries: ("places" | "geometry" | "marker")[] = ["places", "geometry", "marker"];

export interface DeliveryZone {
    radius: number;
    innerRadius: number;
    color: string;
    price: number;
    freeOrderPrice: number;
    name: string;
}

export interface SearchResult {
    zone?: DeliveryZone;
    distance?: number;
    address?: string;
    coords?: { lat: number; lng: number };
}

interface DeliveryMapProps {
    onLocationSelect?: (result: SearchResult) => void;
    cartTotal?: number;
    children?: React.ReactNode;
}

function getCirclePoints(center: google.maps.LatLngLiteral, radiusKm: number, numPoints = 100, clockwise = true) {
    const points: google.maps.LatLngLiteral[] = [];
    const earthRadius = 6371;
    const phase = 2 * Math.PI / numPoints;

    for (let i = 0; i < numPoints; i++) {
        const angle = clockwise ? i * phase : -i * phase;
        const lat1 = (center.lat * Math.PI) / 180;
        const lon1 = (center.lng * Math.PI) / 180;
        const angularDistance = radiusKm / earthRadius;

        const lat2 = Math.asin(
            Math.sin(lat1) * Math.cos(angularDistance) +
            Math.cos(lat1) * Math.sin(angularDistance) * Math.cos(angle)
        );
        const lon2 = lon1 + Math.atan2(
            Math.sin(angle) * Math.sin(angularDistance) * Math.cos(lat1),
            Math.cos(angularDistance) - Math.sin(lat1) * Math.sin(lat2)
        );

        points.push({
            lat: (lat2 * 180) / Math.PI,
            lng: (lon2 * 180) / Math.PI
        });
    }
    points.push(points[0]);
    return points;
}

export default function DeliveryMap({ onLocationSelect, cartTotal, children }: DeliveryMapProps) {
    // 1. Get settings hook
    const { getSetting } = useSettings();
    const { zones } = useDeliveryZones();
    const isMobile = useIsMobile();

    // 2. Initialize state
    const [map, setMap] = useState<google.maps.Map | null>(null);
    const [markerPosition, setMarkerPosition] = useState<google.maps.LatLngLiteral | null>(null);
    const [searchResult, setSearchResult] = useState<SearchResult | null>(null);
    const [mapCenter, setMapCenter] = useState(defaultCenter);
    const [addressValue, setAddressValue] = useState("");

    const { isLoaded } = useJsApiLoader({
        id: 'google-map-script',
        googleMapsApiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || "",
        libraries
    });

    // 3. Update map center from settings
    useEffect(() => {
        const centerSetting = getSetting("zone_center");
        if (centerSetting) {
            try {
                const parsed = JSON.parse(centerSetting.value);
                if (parsed.lat && parsed.lng) {
                    setMapCenter({ lat: parsed.lat, lng: parsed.lng });
                }
            } catch (e) {
                console.error("Failed to parse zone_center setting", e);
            }
        }
    }, [getSetting]);

    const inputRef = useRef<HTMLInputElement>(null);
    const autocompleteRef = useRef<google.maps.places.Autocomplete | null>(null);

    const onLoad = useCallback((map: google.maps.Map) => {
        setMap(map);
    }, []);

    const onUnmount = useCallback(() => {
        setMap(null);
    }, []);

    const checkLocation = useCallback((location: google.maps.LatLngLiteral, address?: string) => {
        const distanceInMeters = google.maps.geometry.spherical.computeDistanceBetween(
            new google.maps.LatLng(mapCenter),
            new google.maps.LatLng(location)
        );
        const distanceInKm = distanceInMeters / 1000;

        const foundZone = zones.find(z => distanceInKm > z.innerRadius && distanceInKm <= z.radius);

        const result: SearchResult = {
            zone: foundZone,
            distance: distanceInKm,
            address: address || "Вибрана точка на мапі",
            coords: { lat: location.lat, lng: location.lng }
        };

        setSearchResult(result);
        setMarkerPosition(location);

        if (onLocationSelect) {
            onLocationSelect(result);
        }

        if (foundZone) {
            toast.success(`Адреса в зоні: ${foundZone.name}`, {
                description: `Вартість доставки: ${foundZone.price}₴ (безкоштовно від ${foundZone.freeOrderPrice}₴)`
            });
        } else {
            toast.error("Адреса поза межами зон доставки", {
                description: "Будь ласка, зв'яжіться з нами для уточнення можливості доставки."
            });
        }
    }, [onLocationSelect, zones, mapCenter]);

    // Initialize Autocomplete
    useEffect(() => {
        if (!isLoaded || !inputRef.current) return;
        if (autocompleteRef.current) return;

        const kharkivBounds = new google.maps.LatLngBounds(
            new google.maps.LatLng(mapCenter.lat - 0.15, mapCenter.lng - 0.2),
            new google.maps.LatLng(mapCenter.lat + 0.15, mapCenter.lng + 0.2)
        );

        const autocomplete = new google.maps.places.Autocomplete(inputRef.current, {
            componentRestrictions: { country: "ua" },
            bounds: kharkivBounds,
            strictBounds: true,
            fields: ['geometry', 'formatted_address', 'name']
        });

        autocomplete.addListener("place_changed", () => {
            const place = autocomplete.getPlace();

            if (!place.geometry || !place.geometry.location) {
                toast.error("Не вдалося знайти деталі для цієї адреси");
                return;
            }

            const location = {
                lat: place.geometry.location.lat(),
                lng: place.geometry.location.lng()
            };

            const address = place.formatted_address || place.name;
            checkLocation(location, address);

            // Update the controlled input value
            if (address) {
                setAddressValue(address);
            }

            if (map) {
                map.panTo(location);
                map.setZoom(15);
            }
        });

        autocompleteRef.current = autocomplete;

    }, [isLoaded, checkLocation, map, mapCenter]);

    const onMapClick = (e: google.maps.MapMouseEvent) => {
        if (e.latLng) {
            const location = {
                lat: e.latLng.lat(),
                lng: e.latLng.lng()
            };

            // Use reverse geocoding to get the address
            const geocoder = new google.maps.Geocoder();
            geocoder.geocode({ location: e.latLng }, (results, status) => {
                let address = "Вибрана точка на мапі";

                if (status === "OK" && results && results[0]) {
                    address = results[0].formatted_address;
                    // Update the controlled input value
                    setAddressValue(address);
                }

                checkLocation(location, address);
            });
        }
    };

    const zonePaths = useMemo(() => {
        return zones.map(zone => {
            const outerPath = getCirclePoints(mapCenter, zone.radius, 120, true);
            const paths = [outerPath];
            if (zone.innerRadius > 0) {
                const innerPath = getCirclePoints(mapCenter, zone.innerRadius, 120, false);
                paths.push(innerPath);
            }
            return paths;
        });
    }, [zones, mapCenter]);

    if (!isLoaded) return <div className="w-full h-[500px] bg-white/5 animate-pulse rounded-2xl" />;

    return (
        <div className="space-y-4">
            <div className="relative z-10">
                <input
                    ref={inputRef}
                    type="text"
                    value={addressValue}
                    onChange={(e) => setAddressValue(e.target.value)}
                    placeholder="Введіть адресу (наприклад: вул. Григорія Сковороди 64)"
                    className="w-full h-12 bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:border-primary transition-colors text-white placeholder:text-gray-500"
                />
            </div>

            {searchResult && (
                <div className={`p-4 rounded-xl border ${searchResult.zone ? 'bg-white/5 border-white/10' : 'bg-red-500/10 border-red-500/20'} animate-in fade-in slide-in-from-top-2`}>
                    {searchResult.zone ? (
                        <div className="space-y-3">
                            <div className="flex justify-between items-start">
                                <div>
                                    <h3 className="font-bold text-lg text-green-400">{searchResult.zone.name}</h3>
                                    <p className="text-sm text-gray-400">Вартість доставки: <span className="text-white font-bold">{searchResult.zone.price} ₴</span></p>
                                </div>
                            </div>

                            {cartTotal !== undefined ? (
                                cartTotal >= searchResult.zone.freeOrderPrice ? (
                                    <div className="p-3 bg-green-500/10 border border-green-500/20 rounded-xl text-green-400 text-sm text-center font-medium">
                                        Безкоштовна доставка
                                    </div>
                                ) : (
                                    <div className="p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-xl text-yellow-500 text-sm text-center">
                                        Додайте товарів ще на <span className="font-bold">{searchResult.zone.freeOrderPrice - cartTotal} ₴</span> для безкоштовної доставки
                                    </div>
                                )
                            ) : (
                                <p className="text-sm text-gray-400 mt-1">Безкоштовно від {searchResult.zone.freeOrderPrice} ₴</p>
                            )}
                        </div>
                    ) : (
                        <div>
                            <h3 className="font-bold text-lg text-red-400">Поза межами зон доставки</h3>
                            <p className="text-sm text-gray-400 mt-1">Ми знаходимось за адресою: вул. Григорія Сковороди 64 (вхід з вул. Багалія)</p>
                        </div>
                    )}
                </div>
            )}

            {children}

            <div className="rounded-2xl overflow-hidden shadow-2xl border border-white/10">
                <GoogleMap
                    mapContainerStyle={containerStyle}
                    center={mapCenter}
                    zoom={11}
                    onLoad={onLoad}
                    onUnmount={onUnmount}
                    onClick={onMapClick}
                    options={{
                        styles: [
                            { featureType: "poi", stylers: [{ visibility: "off" }] },
                            { featureType: "all", elementType: "labels.text.fill", stylers: [{ color: "#9ca3af" }] },
                            { featureType: "all", elementType: "labels.text.stroke", stylers: [{ color: "#242f3e" }] },
                            { featureType: "all", elementType: "geometry", stylers: [{ color: "#242f3e" }] },
                            { featureType: "road", elementType: "geometry", stylers: [{ color: "#38414e" }] },
                            { featureType: "road", elementType: "geometry.stroke", stylers: [{ color: "#212a37" }] },
                            { featureType: "water", elementType: "geometry", stylers: [{ color: "#17263c" }] }
                        ],
                        disableDefaultUI: false,
                        zoomControl: true,
                        streetViewControl: false,
                        mapTypeControl: false,
                        fullscreenControl: true,
                        draggable: true,
                        clickableIcons: false,
                        gestureHandling: isMobile ? 'cooperative' : 'auto'
                    }}
                >
                    {zones.map((zone, index) => (
                        <Polygon
                            key={zone.name}
                            paths={zonePaths[index]}
                            options={{
                                fillColor: zone.color,
                                fillOpacity: 0.35,
                                strokeColor: zone.color,
                                strokeOpacity: 0.8,
                                strokeWeight: 2,
                                clickable: false
                            }}
                        />
                    ))}

                    {markerPosition && (
                        <Marker position={markerPosition} />
                    )}
                </GoogleMap>
            </div>

            <p className="text-sm text-center text-gray-500">
                * Натисніть на мапу, щоб вибрати точку доставки вручну
            </p>
        </div>
    );
}
