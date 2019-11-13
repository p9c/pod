import {readable} from 'svelte/store';

export const config = readable([], function start(set) {
    const interval = setInterval(() => {
        fetch(`http://127.0.0.1:3999/conf`)
            .then(resp => resp.json())
            .then(data => {
                set(data)
            });
    }, 1000);
    return function stop() {
        clearInterval(interval);
    };
});