"use client"

import { GoogleMap, useJsApiLoader, Marker } from "@react-google-maps/api";
import { useState, useEffect } from "react";
import Logo from "@/resources/images/logo.webp";

import { useSettings } from "@/contexts/settings-context";

const containerStyle = {
    width: '100%',
    height: '400px'
};

const defaultCenter = {
    lat: 50.0014656,
    lng: 36.245192
};

// Add 'marker' library
const libraries: ("places" | "geometry" | "marker")[] = ["places", "geometry", "marker"];

export default function WelcomeMap() {
    const { getSetting } = useSettings();
    const [mapCenter, setMapCenter] = useState(defaultCenter);

    const { isLoaded } = useJsApiLoader({
        id: 'google-map-script',
        googleMapsApiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || "",
        libraries
    })

    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        setIsMobile(window.innerWidth < 768);
    }, []);

    // Fetch zone center from settings
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

    if (!isLoaded) return null;

    return (
        <GoogleMap
            mapContainerStyle={containerStyle}
            center={mapCenter}
            zoom={14}
            options={{
                styles: [
                    {
                        featureType: "poi",
                        stylers: [{ visibility: "off" }],
                    },
                    {
                        featureType: "all",
                        elementType: "labels.text.fill",
                        stylers: [{ color: "#9ca3af" }]
                    },
                    {
                        featureType: "all",
                        elementType: "labels.text.stroke",
                        stylers: [{ color: "#242f3e" }]
                    },
                    {
                        featureType: "all",
                        elementType: "geometry",
                        stylers: [{ color: "#242f3e" }]
                    },
                    {
                        featureType: "road",
                        elementType: "geometry",
                        stylers: [{ color: "#38414e" }]
                    },
                    {
                        featureType: "road",
                        elementType: "geometry.stroke",
                        stylers: [{ color: "#212a37" }]
                    },
                    {
                        featureType: "water",
                        elementType: "geometry",
                        stylers: [{ color: "#17263c" }]
                    }
                ],
                draggable: !isMobile,
                streetViewControl: false,
                mapTypeControl: false,
                fullscreenControl: true,
                gestureHandling: isMobile ? 'cooperative' : 'auto'
            }}
        >
            <Marker
                position={mapCenter}
                icon={{
                    url: Logo.src,
                    scaledSize: (typeof google !== 'undefined') ? new google.maps.Size(50, 50) : null
                }}
            />
        </GoogleMap>
    );
}
