import './css/core/reset.css'
import './css/core/zmedia.css'
import './css/core/theme-vars.css'
import './css/common/archive.css'
import './css/common/footer.css'
import './css/common/header.css'
import './css/common/main.css'
import './css/common/not-found.css'
import './css/common/post-entry.css'
import './css/common/post-single.css'
import './css/common/profile-mode.css'
import './css/common/search.css'
import './css/common/terms.css'
import './css/extended/blank.css'
import hljs from "highlight.js";
import 'highlight.js/styles/atom-one-dark.css'

if (window.location.pathname === '/search') {
    import('./js/search')
}

if (window.location.pathname.startsWith('/archives/')) {
    Promise.all([
        import('katex'),
        import('katex/dist/contrib/auto-render'),
        import('katex/dist/katex.min.css'),
    ])
        .then(([katex, autoRenderModule]) => {
            const renderMathInElement = autoRenderModule.default;
            renderMathInElement(document.body, {
                delimiters: [
                    { left: '$$', right: '$$', display: true },
                    { left: '$', right: '$', display: false },
                ],
            });
        })
        .catch((error) => console.error(`Error importing katex: ${error}`));
}

document.addEventListener('DOMContentLoaded', (event) => {
    hljs.highlightAll()
});

