"use client";

import { useEffect, useRef } from "react";

const BASE_CELL_SIZE = 20;

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

/**
 * Determine the best text and cell size to fit the container width.
 * Tries "AUTHFLOW" first, then shorter fallbacks, scaling cell size down if needed.
 */
function pickTextAndCellSize(containerWidth: number): { text: string; cellSize: number } {
    const candidates = ["AUTHFLOW", "AUTH", "AF"];
    const letterWidth = 5;
    const letterSpacing = 2;
    const minPadding = 4; // at least 2 cells padding on each side

    for (const text of candidates) {
        const textWidthCells = text.length * letterWidth + (text.length - 1) * letterSpacing + minPadding;
        // Try the base cell size first, then scale down
        for (let cellSize = BASE_CELL_SIZE; cellSize >= 8; cellSize -= 2) {
            if (textWidthCells * cellSize <= containerWidth) {
                return { text, cellSize };
            }
        }
    }

    // Ultimate fallback — single "A" at smallest size
    return { text: "A", cellSize: 8 };
}

export default function HeroBanner() {
    const heroRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        function generateHeroPixels() {
            const container = heroRef.current;
            if (!container) return;

            const containerWidth = container.clientWidth;
            const containerHeight = container.clientHeight;

            if (containerWidth === 0 || containerHeight === 0) return;

            const { text, cellSize } = pickTextAndCellSize(containerWidth);

            const cols = Math.floor(containerWidth / cellSize);
            const rows = Math.floor(containerHeight / cellSize);

            if (cols < 3 || rows < 5) return;

            container.style.display = "grid";
            container.style.gridTemplateColumns = `repeat(${cols}, ${cellSize}px)`;
            container.style.gridTemplateRows = `repeat(${rows}, ${cellSize}px)`;
            container.innerHTML = "";

            const grid = Array(rows)
                .fill(null)
                .map(() => Array(cols).fill(0));

            const letterWidth = 5;
            const letterSpacing = 2;
            const textWidth =
                text.length * letterWidth + (text.length - 1) * letterSpacing;

            const startX = Math.floor((cols - textWidth) / 2);
            const startY = Math.floor((rows - 5) / 2);

            if (startX >= 0 && startY >= 0) {
                let currentX = startX;
                for (const char of text) {
                    const matrix = pixelFont[char];
                    if (matrix) {
                        for (let r = 0; r < 5; r++) {
                            for (let c = 0; c < 5; c++) {
                                if (
                                    matrix[r][c] === 1 &&
                                    startY + r < rows &&
                                    currentX + c < cols
                                ) {
                                    grid[startY + r][currentX + c] = 1;
                                }
                            }
                        }
                    }
                    currentX += letterWidth + letterSpacing;
                }
            }

            for (let r = 0; r < rows; r++) {
                for (let c = 0; c < cols; c++) {
                    const px = document.createElement("div");
                    if (grid[r][c] === 1) {
                        px.className = "px-fill";
                        const colOffset = c - startX;
                        const delay = colOffset * 15 + Math.random() * 40;
                        px.style.opacity = "0";
                        px.style.transform = "scale(0)";
                        px.style.animation = `pxReveal 0.4s ease-out ${delay}ms forwards`;
                    } else {
                        if (Math.random() > 0.97) {
                            px.className = "px-empty";
                            px.style.opacity = "0";
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
