import { createStore } from 'vuex'

export default createStore({
    state: {
        current: undefined,
        playlist: [],
        recent: [],
        ratings: {}
    },
    mutations: {
        setCurrent(state, c) {
            state.current = c
        },
        setPlaylist(state, p) {
            for (var i=0; i < p.length; i++) {
                state.ratings[p[i].ID] = p[i].Rating
            }
            state.playlist = p
        },
        setRecent(state, r) {
            for (var i=0; i < r.length; i++) {
                state.ratings[r[i].ID] = r[i].Rating
            }
            state.recent = r
        },
        nextTrack(state) {
            // new current track
            var newCurrent = state.playlist.shift()
             
            // move the old one over to the recently played list
            state.recent.unshift(state.current)
            // then trim it
            state.recent.splice(-1, 1)
            state.current = newCurrent
        },
        appendTrack(state, track) {
            state.ratings[track.ID] = track.Rating
            state.playlist.push(track);
        },
        remove(state, idx) {
            state.playlist.splice(idx, 1)
        },
        toTheTop(state, idx) {
            var s = state.playlist[idx]
            state.playlist.splice(idx, 1)
            state.playlist.unshift(s)
        },
        setRating(state, song) {
            state.ratings[song.ID] = song.Rating
        }
    }
})