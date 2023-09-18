# Khatru Invite

A relay based on [Khatru](https://github.com/fiatjaf/khatru) with a invite hierarchy feature.

some notes before running: 
1. change `ws://localhost:3334` in `ui/src/lib/consts.ts` to your relay url endpoint and build the UI
2. configure the relay settings in `.env`
3. manually add someone to the `whitelist.json` file, like this: `[{"pk":"07adfda9c5adc80881bb2a5220f6e3181e0c043b90fa115c4f183464022968e6","invited_by":""}]`