package com.datatheorem.logproc;

import lombok.Getter;
import lombok.Setter;
import java.util.Map;

public class VisitorEntry {
    @Getter @Setter private String Timestamp;
    @Getter @Setter private String Source;
    @Getter @Setter private String FName;
    @Getter @Setter private String LName;
    @Getter @Setter private String City;
    @Getter @Setter private String State;
    @Getter @Setter private String Country;
    @Getter @Setter private String Message;

    public VisitorEntry(Map<String, String> values) {
        // Time,Source,First Name,Last Name,City,State,Country,Message
        Timestamp = values.get("Time");
        Source = values.get("Source");
        FName = values.get("First Name");
        LName = values.get("Last Name");
        City = values.get("City");
        State = values.get("State");
        Country = values.get("Country");
        Message = values.get("Message");
    }

    public String toString() {
        return Timestamp + " (" + Source + "): " + FName + " " + LName + " from " + City + ", " + State + " (" + Country + ") wrote: \"" + Message + "\"";
    }
}
