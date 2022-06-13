package com.example.welcome;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;

@ControllerAdvice
class GiftCardNotFoundAdvice {

    @ResponseBody
    @ExceptionHandler(GiftCardNotFoundException.class)
    @ResponseStatus(HttpStatus.NOT_FOUND)
    String giftCardNotFoundHandler(GiftCardNotFoundException ex) {
        return ex.getMessage();
    }
}
