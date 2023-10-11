import './style.css';
import { Relay, relayInit } from 'nostr-tools';
import { loadActions } from './actions';

export let relayUrl = '';
export let relay: Relay;

window.onload = function () {
  setTimeout(() => {
    // get relay websocket url
    relayUrl = '';
    var loc = window.location;
    if (loc.protocol === "https:") {
      relayUrl = "wss:";
    } else {
      relayUrl = "ws:";
    }
    relayUrl += "//" + loc.host;

    // init relay connection
    relay = relayInit(relayUrl);
    relay.on('connect', () => {
      console.log(`connected to ${relay.url}`);
    });
    relay.on('error', () => {
      console.log(`failed to connect to ${relay.url}`);
    });
    relay.connect();

    // 
    loadActions();
  }, 5);
};

