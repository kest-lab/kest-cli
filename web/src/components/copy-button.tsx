'use client';

import { Button } from '@/components/ui/button';

interface CopyButtonProps {
    text: string;
    variant?: 'default' | 'ghost' | 'outline';
    size?: 'default' | 'sm' | 'lg' | 'icon';
    className?: string;
    children?: React.ReactNode;
}

export function CopyButton({ text, variant = 'ghost', size = 'sm', className, children }: CopyButtonProps) {
    const handleCopy = () => {
        navigator.clipboard.writeText(text);
    };

    return (
        <Button
            variant={variant}
            size={size}
            className={className}
            onClick={handleCopy}
        >
            {children || 'Copy'}
        </Button>
    );
}
