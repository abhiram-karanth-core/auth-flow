"use client";

import { useEffect, useRef, useState } from "react";

function getColumnLabel(index: number): string {
  let label = "";
  let i = index + 1;
  while (i > 0) {
    const rem = (i - 1) % 26;
    label = String.fromCharCode(65 + rem) + label;
    i = Math.floor((i - 1) / 26);
  }
  return label;
}

const CELL_SIZE = 20;

export default function GridBackground() {
  const rulerXRef = useRef<HTMLDivElement>(null);
  const rulerYRef = useRef<HTMLDivElement>(null);
  const cursorBoxRef = useRef<HTMLDivElement>(null);
  const [coords, setCoords] = useState({ x: "A", y: 1 });

  useEffect(() => {
    function generateRulers() {
      const width = window.innerWidth;
      const cols = Math.ceil(width / CELL_SIZE);

      if (rulerXRef.current) {
        rulerXRef.current.innerHTML = "";
        for (let i = 0; i < cols; i++) {
          const span = document.createElement("span");
          span.textContent = getColumnLabel(i);
          rulerXRef.current.appendChild(span);
        }
      }

      const rows = Math.ceil(document.body.scrollHeight / CELL_SIZE);
      if (rulerYRef.current) {
        rulerYRef.current.innerHTML = "";
        for (let i = 0; i < rows; i++) {
          const span = document.createElement("span");
          span.textContent = String(i + 1);
          rulerYRef.current.appendChild(span);
        }
      }
    }

    function handleMouseMove(e: MouseEvent) {
      const x = Math.floor((e.clientX - 30) / CELL_SIZE);
      const y = Math.floor((e.clientY - 20) / CELL_SIZE);
      const visualX = x * CELL_SIZE + 30;
      const visualY = y * CELL_SIZE + 20;

      if (e.clientX > 30 && e.clientY > 20) {
        if (cursorBoxRef.current) {
          cursorBoxRef.current.style.display = "block";
          cursorBoxRef.current.style.transform = `translate(${visualX}px, ${visualY}px)`;
        }
        setCoords({ x: getColumnLabel(x), y: y + 1 });
      } else {
        if (cursorBoxRef.current) {
          cursorBoxRef.current.style.display = "none";
        }
      }
    }

    generateRulers();
    window.addEventListener("resize", generateRulers);
    document.addEventListener("mousemove", handleMouseMove);

    return () => {
      window.removeEventListener("resize", generateRulers);
      document.removeEventListener("mousemove", handleMouseMove);
    };
  }, []);

  return (
    <>
      <div className="corner-piece" />
      <div className="ruler-x" ref={rulerXRef} />
      <div className="ruler-y" ref={rulerYRef} />
      <div className="grid-layer" />
      <div className="active-cell-indicator" ref={cursorBoxRef} />
      <div className="coordinates-display">
        X: {coords.x} | Y: {coords.y}
      </div>
    </>
  );
}
