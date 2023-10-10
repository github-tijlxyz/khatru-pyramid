import { NDKEvent, NDKRelay, NDKRelaySet } from '@nostr-dev-kit/ndk';
import { nip19 } from 'nostr-tools';
import { ndk, relayUrl } from './main';

// Handle Buttons
export function loadActions() {
    document.querySelectorAll(".removefromrelay").forEach((btn) => {
        btn.addEventListener('click', () => {
            const arg = btn.getAttribute('data-actionarg');
            if (arg) {
                const jsonArg = JSON.parse(arg);
                if (jsonArg) {
                    removeFromRelay(jsonArg)
                }
            };
        })
    })
    document.querySelectorAll(".inviteuser").forEach((btn) => {
        btn.addEventListener('click', () => {
            inviteUser();
        })
    })

    async function inviteUser() {
        // only publish to the relay in question, dont know why this needs so much code
        let specificRelay = new Set<NDKRelay>
        specificRelay.add(new NDKRelay(relayUrl))
        const relaySet = new NDKRelaySet(specificRelay, ndk);
        relaySet.relays.forEach(async (relay) => {
            await relay.connect().catch((err) => {
                console.log("RELAY CONNECT ERROR");
                console.error(err);
            });
            relay.on("connect", () => {
                console.log(relay.url, "connected");
            });
        });

        var input = (<HTMLInputElement>document.getElementById('inviteuser-input')).value;
        let p = input;
        if (input.startsWith("npub")) {
            let h = nip19.decode(input).data;
            if (typeof h == 'string') {
                p = h;
            }
        }
        const event = new NDKEvent(ndk)
        event.kind = 20201;
        event.tags = [['p', p]];
        event.content = "";
        const relays = await event.publish(relaySet);
        relays.forEach(() => {
            setTimeout(() => {
                window.location.href = ""
            }, 128)
        });
    }

    async function removeFromRelay(tags: string[][]) {
        // only publish to the relay in question, dont know why this needs so much code
        let specificRelay = new Set<NDKRelay>
        specificRelay.add(new NDKRelay(relayUrl))
        const relaySet = new NDKRelaySet(specificRelay, ndk);
        relaySet.relays.forEach(async (relay) => {
            await relay.connect().catch((err) => {
                console.log("RELAY CONNECT ERROR");
                console.error(err);
            });
            relay.on("connect", () => {
                console.log(relay.url, "connected");
            });
        });

        const event = new NDKEvent(ndk)
        event.kind = 20202;
        event.tags = tags;
        event.content = "";
        const relays = await event.publish();
        relays.forEach(() => {
            setTimeout(() => {
                window.location.href = ""
            }, 128)
        });
    }
}