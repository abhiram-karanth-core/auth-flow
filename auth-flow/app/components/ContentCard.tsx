import { ReactNode } from "react";

interface ContentCardProps {
    children: ReactNode;
    id?: string;
    className?: string;
    style?: React.CSSProperties;
}

export default function ContentCard({ children, id, className = "", style }: ContentCardProps) {
    return (
        <div className={`content-card ${className}`} id={id} style={style}>
            {children}
        </div>
    );
}

interface TitleBlockProps {
    meta: string;
    title: string;
}

export function TitleBlock({ meta, title }: TitleBlockProps) {
    return (
        <div className="title-block">
            <span className="meta-data">{meta}</span>
            <h2>{title}</h2>
        </div>
    );
}
