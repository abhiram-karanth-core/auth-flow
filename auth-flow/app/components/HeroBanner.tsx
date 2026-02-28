"use client";

import { useEffect, useRef } from "react";

const CELL_SIZE = 20;

const pixelFont: Record<string, number[][]> = {
    A: [
        [0, 1, 1, 1, 0],
        [1, 0, 0, 0, 1],
        [1, 1, 1, 1, 1],
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
    ],
    U: [
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [0, 1, 1, 1, 0],
    ],
    T: [
        [1, 1, 1, 1, 1],
        [0, 0, 1, 0, 0],
        [0, 0, 1, 0, 0],
        [0, 0, 1, 0, 0],
        [0, 0, 1, 0, 0],
    ],
    H: [
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [1, 1, 1, 1, 1],
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
    ],
    F: [
        [1, 1, 1, 1, 1],
        [1, 0, 0, 0, 0],
        [1, 1, 1, 1, 0],
        [1, 0, 0, 0, 0],
        [1, 0, 0, 0, 0],
    ],
    L: [
        [1, 0, 0, 0, 0],
        [1, 0, 0, 0, 0],
        [1, 0, 0, 0, 0],
        [1, 0, 0, 0, 0],
        [1, 1, 1, 1, 1],
    ],
    O: [
        [0, 1, 1, 1, 0],
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [0, 1, 1, 1, 0],
    ],
    W: [
        [1, 0, 0, 0, 1],
        [1, 0, 0, 0, 1],
        [1, 0, 1, 0, 1],
        [1, 0, 1, 0, 1],
        [0, 1, 0, 1, 0],
    ],
};

// Inject keyframes once
if (typeof document !== "undefined") {
    const styleId = "hero-pixel-anim";
    if (!document.getElementById(styleId)) {
        const style = document.createElement("style");
        style.id = styleId;
        style.textContent = `
            @keyframes pxReveal {
                0%   { opacity: 0; transform: scale(0); }
                60%  { opacity: 1; transform: scale(1.15); }
                100% { opacity: 1; transform: scale(1); }
            }
            @keyframes pxFadeIn {
                0%   { opacity: 0; }
                100% { opacity: 0.3; }
            }
        `;
        document.head.appendChild(style);
    }
}

export default function HeroBanner() {
    const heroRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        function generateHeroPixels() {
            const container = heroRef.current;
            if (!container) return;

            const containerWidth = container.clientWidth;
            const containerHeight = container.clientHeight;
            const cols = Math.floor(containerWidth / CELL_SIZE);
            const rows = Math.floor(containerHeight / CELL_SIZE);

            container.style.display = "grid";
            container.style.gridTemplateColumns = `repeat(${cols}, ${CELL_SIZE}px)`;
            container.style.gridTemplateRows = `repeat(${rows}, ${CELL_SIZE}px)`;
            container.innerHTML = "";

            const grid = Array(rows)
                .fill(null)
                .map(() => Array(cols).fill(0));

            const text = "AUTHFLOW";
            const letterWidth = 5;
            const letterSpacing = 2;
            const textWidth =
                text.length * letterWidth + (text.length - 1) * letterSpacing;

            const startX = Math.floor((cols - textWidth) / 2);
            const startY = Math.floor((rows - 5) / 2);

            if (startX > 0 && startY > 0) {
                let currentX = startX;
                for (const char of text) {
                    const matrix = pixelFont[char];
                    if (matrix) {
                        for (let r = 0; r < 5; r++) {
                            for (let c = 0; c < 5; c++) {
                                if (matrix[r][c] === 1) {
                                    grid[startY + r][currentX + c] = 1;
                                }
                            }
                        }
                    }
                    currentX += letterWidth + letterSpacing;
                }
            }

            // Track column index of text pixels for staggered delay
            let textPixelIndex = 0;

            for (let r = 0; r < rows; r++) {
                for (let c = 0; c < cols; c++) {
                    const px = document.createElement("div");
                    if (grid[r][c] === 1) {
                        px.className = "px-fill";
                        // Stagger based on column position relative to text start
                        const colOffset = c - startX;
                        const delay = colOffset * 15 + Math.random() * 40;
                        px.style.opacity = "0";
                        px.style.transform = "scale(0)";
                        px.style.animation = `pxReveal 0.4s ease-out ${delay}ms forwards`;
                        textPixelIndex++;
                    } else {
                        if (Math.random() > 0.97) {
                            px.className = "px-empty";
                            px.style.opacity = "0";
                            // Scatter pixels appear after text finishes
                            const scatterDelay = 800 + Math.random() * 600;
                            px.style.animation = `pxFadeIn 0.6s ease-out ${scatterDelay}ms forwards`;
                        }
                    }
                    container.appendChild(px);
                }
            }
        }

        generateHeroPixels();
        window.addEventListener("resize", generateHeroPixels);

        return () => {
            window.removeEventListener("resize", generateHeroPixels);
        };
    }, []);

    return <div className="hero-banner" ref={heroRef} />;
}
