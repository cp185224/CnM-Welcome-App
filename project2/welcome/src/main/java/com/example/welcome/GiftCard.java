package com.example.welcome;

import org.apache.tomcat.jni.Local;

import java.util.Objects;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import java.util.Date;

@Entity
public class GiftCard {
    private @Id @GeneratedValue Long id;
    private double balance;
    private Date timeCreated;

    public GiftCard() {}

    public GiftCard(double balance) {

        this.balance = balance;
        this.timeCreated = new Date();
    }

    public GiftCard(double balance, Date timeCreated) {

        this.balance = balance;
        this.timeCreated = timeCreated;
    }

    public Long getId() {
        return this.id;
    }

    public double getBalance() {
        return this.balance;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public void setBalance(double balance) {
        this.balance = balance;
    }

    public Date getTimeCreated() {
        return timeCreated;
    }

    public void setTimeCreated(Date timeCreated) {
        this.timeCreated = timeCreated;
    }

    @Override
    public boolean equals(Object o) {

        if (this == o)
            return true;
        if (!(o instanceof GiftCard))
            return false;
        GiftCard giftCard = (GiftCard) o;
        return Objects.equals(this.id, giftCard.id) && Objects.equals(this.balance, giftCard.balance);
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.id, this.balance);
    }

    @Override
    public String toString() {
        return "GiftCard{" + "id=" + this.id + ", balance='" + this.balance + '\'' + '}';
    }
}

