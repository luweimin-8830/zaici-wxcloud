export function notFoundHandler(req, res) {
    if (req.path.startsWith('/api')) {
        return res.status(404).json({ code: 404, message: 'Not Found' });
    }
    res.status(404).send('404 Not Found');
}

export function errorHandler(err, req, res) {
    console.error(err);
    if (req.path.startsWith('/api')) {
        res.status(500).json({ code: 500, message: err.message || 'Internal Server Error' });
    } else {
        res.status(500).send('Server Error');
    }
}