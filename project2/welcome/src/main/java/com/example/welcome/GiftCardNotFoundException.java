package com.example.welcome;

public class GiftCardNotFoundException extends RuntimeException {
    GiftCardNotFoundException(Long id) {
        super("Could not find gift card " + id);
    }
}
