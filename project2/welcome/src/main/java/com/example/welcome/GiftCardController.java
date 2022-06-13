package com.example.welcome;

import java.util.Date;
import java.util.List;

import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class GiftCardController {

    private final GiftCardRepository repository;
    private final long YEARINMS = 120000L;
//    private final long YEARINMS = 31556926000L;

    public GiftCardController(GiftCardRepository repository) {
        this.repository = repository;
    }


    private boolean isExpired(GiftCard giftCard) {
        Long ageInMs = (new Date()).getTime() - giftCard.getTimeCreated().getTime();
        return ageInMs > YEARINMS;
    }

    // Aggregate root
    // tag::get-aggregate-root[]
    @GetMapping("/giftcards")
    List<GiftCard> all() {
        List<GiftCard> giftCards = repository.findAll();
        for (GiftCard giftCard : giftCards) {
            if (isExpired(giftCard)) {
                long id = giftCard.getId();
                deleteGiftCard(id);
            }
        }
        return repository.findAll();
    }
    // end::get-aggregate-root[]

    @PostMapping("/giftcards")
    GiftCard newGiftCard(@RequestBody GiftCard newGiftCard) {
        if (newGiftCard.getTimeCreated() == null) {
            newGiftCard.setTimeCreated(new Date());
        }
        return repository.save(newGiftCard);
    }

    // Single item

    @GetMapping("/giftcards/{id}")
    GiftCard one(@PathVariable Long id) {
        GiftCard oneGiftCard = repository.findById(id)
                .orElseThrow(() -> new GiftCardNotFoundException(id));
        if (isExpired(oneGiftCard)) {
            deleteGiftCard(id);
            throw new GiftCardNotFoundException(id);
        }
        return oneGiftCard;
    }

    @PutMapping("/giftcards/{id}")
    GiftCard replaceGiftCard(@RequestBody GiftCard newGiftCard, @PathVariable Long id) {
        if (newGiftCard.getTimeCreated() == null) {
            newGiftCard.setTimeCreated(new Date());
        }
        return repository.findById(id)
                .map(giftCard -> {
                    giftCard.setBalance(newGiftCard.getBalance());
                    return repository.save(giftCard);
                })
                .orElseGet(() -> {
                    newGiftCard.setId(id);
                    return repository.save(newGiftCard);
                });
    }

    @DeleteMapping("/giftcards/{id}")
    void deleteGiftCard(@PathVariable Long id) {
        repository.deleteById(id);
    }
    // Return 204
}
