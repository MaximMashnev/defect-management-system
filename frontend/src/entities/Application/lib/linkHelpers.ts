import { ApplicationEntity } from "../model/types";

export function linkGeneration (app: ApplicationEntity): string {
    const emailText = app.email.split('@')[0].slice(1, 3);
    return `a${app.id}b${app.building_id}e${emailText}p`;
}

export function linkDecrypting (link: string) {
    try {
        const match = link.match(/^a(\d+)b/);
        if (match && match[1]) {
            return Number(match[1]);
        }
        return null;
    } catch {
        return null;
    }
}