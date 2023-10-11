import { nip19 } from 'nostr-tools';
import { relay } from './main';

export function loadActions() {

    // Handle removefromrelay buttons
    document.querySelectorAll(".removefromrelay").forEach((btn) => {
        btn.addEventListener('click', () => {
            const arg = btn.getAttribute('data-actionarg');
            if (arg) {
                const jsonArg = JSON.parse(arg);
                if (jsonArg) removeFromRelay(jsonArg);
            };
        });
    });

    //handle inviteuser buttons
    document.querySelectorAll(".inviteuser").forEach((btn) => {
        btn.addEventListener('click', () => {
            inviteUser();
        });
    });

    async function inviteUser() {

        var input = (<HTMLInputElement>document.getElementById('inviteuser-input')).value;
        let p = input;
        if (input.startsWith("npub")) {
            let h = nip19.decode(input).data;
            if (typeof h == 'string') {
                p = h;
            };
        };

        let event = {
            kind: 20201,
            created_at: Math.floor(Date.now() / 1000),
            content: '',
            tags: [['p', p]],
        };

        // @ts-ignore for window.nostr
        const signedEvent = await window.nostr.signEvent(event);

        await relay.publish(signedEvent).then(() => {
            setTimeout(() => window.location.href = "", 300);
        });
    };

    async function removeFromRelay(tags: string[][]) {
        let event = {
            kind: 20202,
            created_at: Math.floor(Date.now() / 1000),
            content: '',
            tags,
        };

        // @ts-ignore for window.nostr
        const signedEvent = await window.nostr.signEvent(event);

        await relay.publish(signedEvent).then(() => {
            setTimeout(() => window.location.href = "", 300);
        });
    };
};