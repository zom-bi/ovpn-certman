$.ajaxSetup({
    // Initialize headers from the csrf-token meta tag.
    headers: {
        'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')
    }
});