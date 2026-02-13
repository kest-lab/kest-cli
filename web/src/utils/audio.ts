import { env } from '@/config/env';

/**
 * Audio Utility
 * Provides methods to play sound effects for UI feedback
 */

const playSound = (path: string) => {
    if (typeof window === 'undefined') return;
    const audio = new Audio(path);
    audio.play().catch(err => {
        // Browsers often block audio play without user interaction
        console.warn('Audio playback failed:', err);
    });
};

const speak = async (text: string) => {
    if (typeof window === 'undefined') return;

    try {
        const response = await fetch('/api/tts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ text })
        });

        if (!response.ok) {
            throw new Error(`TTS request failed: ${response.statusText}`);
        }

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const audio = new Audio(url);
        await audio.play();
    } catch (err) {
        console.error('TTS playback failed:', err);
        // Fallback to simple sound if TTS fails
        playSound('/audio/success.mp3');
    }
};

export const audioUtils = {
    success: (text?: string) => {
        if (text) {
            speak(text);
        } else {
            playSound('/audio/success.mp3');
        }
    },
    error: (text?: string) => {
        if (text) {
            speak(text);
        } else {
            playSound('/audio/error.mp3');
        }
    },
    speak,
};
