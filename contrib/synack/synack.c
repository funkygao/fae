/*
**                                Copyright 2000
**                                      by
**                         The Board of Trustees of the
**                       Leland Stanford Junior University.
**                              All rights reserved.
**
**
**                               Disclaimer Notice
**
**   The items furnished herewith were developed under the sponsorship
**   of the U.S. Government.  Neither the U.S., nor the U.S. D.O.E., nor
**   the Leland Stanford Junior University, nor their employees, makes any 
**   warranty, express or implied, or assumes any liability or 
**   responsibility for accuracy, completeness or usefulness of any
**   information, apparatus, product or process disclosed, or represents
**   that its use will not infringe privately-owned rights.  Mention of any 
**   product, its manufacturer, or suppliers shall not, nor is it intended to, 
**   imply approval, disapproval, or fitness for any particular use.  The U.S. 
**   and the University at all times retain the right to use and disseminate the
**   furnished items for any purpose whatsoever.
**
**   Notice 91 02 01
**
**   Work supported by the U.S. Department of Energy under contract
**   DE-AC03-76SF00515.
*/

/* 
**   SYNACK tool developed by Kishan Jayaraman
*/ 

#include <stdio.h>
#include <string.h> 
#include <ctype.h>
#include <signal.h>
#include <sys/types.h> 
#include <sys/socket.h> 
#include <stdlib.h> 
#include <errno.h> 
#include <netdb.h> 
#include <netinet/in.h> 
#include <sys/time.h>
#include <math.h>
#include <fcntl.h>
#include <sysexits.h>
#include <pthread.h>

#define MYPORT 5000
#define HISPORT 22
#define MYIP "doris.slac.stanford.edu"
#define HISIP "doris.slac.stanford.edu"
#define BACKLOG 1
#define MAXSECS         2146    /* 2147,483,647 usec */
#define VERY_LONG       ((time_t)MAXSECS*1000000)
#define MAXLINE 1
#define TRUE 1
#define FALSE 0
#define MAX_CON 100000


struct sockaddr_in my_addr;
struct sockaddr_in his_addr;
struct hostent *myadr, *hisadr;
struct in_addr inaddr;
struct addrinfo *webadr;

time_t tvsub();
char * tvprint();
void err();
void record_stats();
void statistics();
void print_stats();
void printstats();
void callstats();
double xsqrt();
int isnumeric();
void rterror();
void err_print();
void sock_init();
void sock_close();
struct timeval* get_time();
int connect_start();
int select_start();
int select_condition();
void connection_tcp();
void *call_connection();
int max_threads();
void finish();
void exit_interrupt();
void interval_between_synacks();

int percent1 = 25;              /* 1st percentile */
int percent2 = 75;       	/* 2nd percentile */
int hisport = 22;		/* port number of server */ 

int counter=0;
int sockfd[MAX_CON], flags[MAX_CON];
char *ipname, *ipnum;
fd_set rset[MAX_CON], wset[MAX_CON];
struct timeval tvalfresh;
int secs = 0;			/* secs for timeout */
long int millisecs=0;		/* millisecondes for timeout */
int threads;
pthread_t tid[MAX_CON];
int key = 0;  

/* Statistics Structure */ 

  typedef struct _stats {
        time_t rttmin;                  /* minimum round-trip time */
        time_t rttmax;                  /* maximum round-trip time */
        double rttsum;                  /* sum of recorded rtt values */
        double rttssq;                  /* sum of squared  rtt values */
        time_t rttmed;                  /* median of rtt values*/
        time_t rttiqr;                  /* interquartile range of rtt values */
        time_t rttpct1;                 /* rtt value for the 1st %ile */
        time_t rttpct2;                 /* rtt value for the 2nd %ile */
        int rcvd;                       /* no. of successful connections */
	int seq_no;			/* sequence number */ 	
        time_t rttarray[1000000];       /* array of rtt values */
  }stats;

  stats synack = { VERY_LONG, 0 , 0, 0, 0, 0, 0, 0, 0, 0, 0};

typedef struct timer {
	int sec;
	int usec;
}timer;


extern char *optarg;
extern int optind;

char Usage[] = "\
Usage: synack [-options] host\n\
Common options:\n\
-p ##   port number to send to (default 22)\n\
-k ##   no. of connections to be made\n\
-i ##   Time interval between connections in secs (default 1 sec)\n\
-u ##   Time interval between connections in microsecs (minimum 20000 microsecs)\n\
-z ##   Percentile 1 (default 25)\n\
-Z ##   Percentile 2 (default 75)\n\
-S ##	Timout in secs (default 1 Sec)\n\
-s ##	Timeout in millisecs \n\
";

int 
main(argc,argv)
int argc;
char **argv;
{

/* ********************************************************************** */

  int sin_size,c;
  int nconnect,nselect,nselectval,nread,nweb,nset;
  int iscount = FALSE;
  long int count = 1;
  long int i,j;
  int sig = 0;
  int time = 0;
  unsigned long gw=0;

	
  char myip[] = MYIP;
  char *hisip ;
  char *wmaddr = "wmatthews=home-1.stanford.edu";
  char *wmport = "80";
  char line[MAXLINE +1];

  struct timeval *sendtime;  
  struct timeval *recvtime;
  struct timeval tv1;
  struct timeval tv2;
  struct timeval tval;
 
  timer interval;
 
  time_t rtt = 0;


  /*  Command line options */

  long int length = 0;		/* default value */
  int interval_sec = 0;		/* time interval between connections in secs */	
  int interval_usec = 0;	/* time interval between connections in microsecs */


/* **************************************************************************** */


  	if (argc < 1) goto usage;
  
  	while ((c = getopt(argc, argv, "p:k:i:u:z:Z:s:S:")) != -1) {
    		switch (c) { 

    		case 'p':
      			hisport = atoi(optarg);
    		break;

    		case 'k':
      			length = atoi(optarg);
    		break;

    		case 'i':
      			interval_sec = atoi(optarg);
    		break;

    		case 'u':
      			interval_usec = atoi(optarg);
    		break;	

    		case 'z':
      			percent1 = atoi(optarg);
   	 	break;
  
    		case 'Z':
      			percent2 = atoi(optarg);
    		break;
  
    		case 'S':
      			secs = atoi(optarg);
    		break;

    		case 's':
      			millisecs = atoi(optarg);
    		break;
    
    		default:
      			goto usage;
   		}
  	}
 
	if ( (interval_sec == 0) && (interval_usec == 0) ) {
		interval.sec = 1;
		interval.usec = 0;
	}
	else {
		interval.sec = interval_sec;
		interval.usec = interval_usec;
	}

	if ( (interval_sec == 0) && (interval_usec < 20000) )
		interval.usec = 20000;

  	if (length > 0) 
    		iscount = TRUE;

  	if (secs == 0 && millisecs == 0) 
		secs = 10;

  	tvalfresh.tv_sec = secs;
  	tvalfresh.tv_usec = millisecs * 1000;

  	if (optind == argc)
    		goto usage;
  	hisip = argv[optind];

  	myadr = gethostbyname(myip);
  	hisadr = gethostbyname(hisip);

  	my_addr.sin_family = AF_INET;         
  	my_addr.sin_port = htons(MYPORT);    
  	my_addr.sin_addr.s_addr = htonl(INADDR_ANY);
  	myadr = gethostbyaddr((char *) &my_addr.sin_addr.s_addr, sizeof(struct in_addr*), AF_INET);
  	bzero(&(my_addr.sin_zero), 8);          

  	if(!hisadr) {
     		rterror("Unknown host %s", hisip);
      		exit(EX_NOHOST);
  	} 


  	his_addr.sin_family = AF_INET;
  	his_addr.sin_port = htons(hisport);
  	his_addr.sin_addr.s_addr = inet_addr(inet_ntoa(*((struct in_addr*)hisadr->h_addr)));
  
  	if (isnumeric(hisip)) {
		hisadr->h_name = NULL;
		hisadr = gethostbyaddr((char *) &his_addr.sin_addr.s_addr, sizeof(struct in_addr*), AF_INET);
  	}

  	bzero(&(his_addr.sin_zero), 8); 
  
  	if (isnumeric(hisip) && !hisadr) {
		printf("\nAddress Information of Server not listed in Domain\n");  
		hisadr = gethostbyname(hisip);
  	}

  	if (!isnumeric(hisadr->h_name))
		ipname = hisadr->h_name;
  	else
		ipname = NULL; 

  	ipnum = (char *)inet_ntoa(*((struct in_addr *)hisadr->h_addr));

  	if (iscount) 
		if (ipname != NULL)
			printf("\nSYN-ACK to %s (%s), %d Packets\n\n", ipname, ipnum, length);
		else  
			printf("\nSYN-ACK to %s , %d Packets\n\n", ipnum, length);
  	else 
		if (ipname != NULL)
			printf("\nSYN-ACK to %s (%s)\n\n", ipname, ipnum);   
		else 
			printf("\nSYN-ACK to %s \n\n", ipnum);
	

/*  
**  BEGINNNING OF CONNECT LOOP 
**  -----------------------------
*/
  
  	while ( (!sig) && (((iscount) && (count <= length)) || (!iscount)) ) {  

	
		threads = max_threads(secs,interval);

		i = pthread_create(&tid[counter], NULL, call_connection, NULL);

		j = pthread_detach(tid[counter]);

		if (i == EAGAIN) {
			printf("\n\nOOPS....LIMIT REACHED ON MAX. NO OF THREADS :( \n\n");
			goto out_of_loop;
		}

		if (counter == threads) 
			counter = 0;
		else 
			counter++;	  

 		interval_between_synacks(interval);

  		count++;

		synack.seq_no++;

  		signal(SIGINT, exit_interrupt); 

		if (key == 1) {
			finish();
		}

  	}

/*
** END OF CONNECT LOOP 
** -------------------
*/

  out_of_loop:

  	printf("\nWaiting for outstanding packets (if any)..........\n\n");
  
  	for (i=0; i <= threads; i++) 
		pthread_join(tid[i], NULL);	

  	sleep(1);

  	print_stats(&synack); 
  	exit(0);

  	usage:
   		fprintf(stderr,Usage);
  	exit(1);

}


/*
** END OF MAIN() 
** -------------
*/




/*
** FUNCTIONS 
** ---------
*/


/* 
** ERR -- Prints standard error 
** ----------------------------
*/

void  
err(s)
     char *s;
{ 
  	fprintf(stderr,"synack%: ");
  	perror(s);
  	fprintf(stderr,"errno=%d\n",errno);
  	exit(1);
}


/*
** TVSUB -- Subtract two timeval structs
** -------------------------------------
**
**      Returns:
**              Time difference in micro(!)seconds.
**
**      Side effects:
**              The difference of the two timeval structs
**              is stored back into the first.
**
**      This implementation assumes that time_t is a signed value.
**      On 32-bit machines this limits the range to ~35 minutes.
**      That seems sufficient for most practical purposes.
**      Note that tv_sec is an *un*signed entity on some platforms.
*/
  
time_t
tvsub(t2, t1)
struct timeval *t2;                     
struct timeval *t1;                     
{ 
        register time_t usec;

        t2->tv_usec -= t1->tv_usec;
        while (t2->tv_usec < 0)
        {
                t2->tv_usec += 1000000;
                if (t2->tv_sec != 0)
                        t2->tv_sec--;
                else
                        t2->tv_usec = 0;
        }
  
        if (t2->tv_sec < t1->tv_sec)
        {
                t2->tv_sec = 0;
                t2->tv_usec = 0;
        }
        else
                t2->tv_sec -= t1->tv_sec;
        
        if (t2->tv_sec > MAXSECS)  
        {
                t2->tv_sec = MAXSECS;
                t2->tv_usec = 0;
        }
                        
        usec = t2->tv_sec*1000000 + t2->tv_usec;
        return(usec);
}


/*
** TVPRINT -- Convert time value to ascii string
** ---------------------------------------------
**              
**      Returns:
**              Pointer to string in static storage.
**
**      Output is in variable format, depending on the value.
**      This avoids printing of non-significant digits.
*/

char *
tvprint(usec)           
time_t usec;                            /* value to convert */
{                       
        static char buf[30];            /* sufficient for 64-bit values */
        time_t uval;

        uval = usec;
        (void) sprintf(buf, "%ld.%3.3ld", uval/1000, uval%1000);
                
        return(buf);
}


/* 
** RECORD_STATS -- Records RTT Statistics (Min, Mean, Max, Standard Deviation)
** ---------------------------------------------------------------------------------
*/ 

void
record_stats(sp, rtt)
stats *sp;                              /* statistics buffer */
time_t rtt;                             /* round-trip time */
{
        if (rtt < sp->rttmin)
                sp->rttmin = rtt;
         
        if (rtt > sp->rttmax)
                sp->rttmax = rtt;

        sp->rttarray[sp->rcvd-1] = rtt;
        sp->rttsum += (double)rtt;
        sp->rttssq += (double)rtt * (double)rtt;
}               


/*
** STATISTICS -- Records RTT Statistics (Median, Interquartile Range, Percentiles)
** -------------------------------------------------------------------------------
*/
 
void
statistics(sp,percent1,percent2)
stats *sp;
int percent1,percent2;
{
        int loop1,loop2,q1,q3;
        time_t temp;
                        
        for(loop1=0; loop1 <= sp->rcvd-2; loop1++)
        for(loop2=loop1+1; loop2 <= sp->rcvd-1; loop2++)
        {
                if(sp->rttarray[loop1] > sp->rttarray[loop2])
                {
                        temp = sp->rttarray[loop1];
                        sp->rttarray[loop1] = sp->rttarray[loop2];
                        sp->rttarray[loop2] = temp;
                }
        } 
        if((sp->rcvd%2)==0)
                sp->rttmed = ((sp->rttarray[sp->rcvd/2-1])+(sp->rttarray[sp->rcvd/2]))/2;
        else
                sp->rttmed = sp->rttarray[sp->rcvd/2];
        if((sp->rcvd%4) == 0)
                {
                        q1 = sp->rcvd/4 - 1;
                        q3 = (3 * (sp->rcvd/4)) - 1 ;
                }
        else
                {
                        q1 = sp->rcvd/4;
                        q3 = (int)(3 * ((double)sp->rcvd/4.0));
                }
        sp->rttiqr = sp->rttarray[q3] - sp->rttarray[q1];
                
        if ((sp->rcvd * percent1) >= 100)
                if (((sp->rcvd * percent1) % 100) == 0)
                        sp->rttpct1 = sp->rttarray[(sp->rcvd * percent1 / 100) - 1];
                else
                        sp->rttpct1 = sp->rttarray[sp->rcvd * percent1 / 100];
        else
                sp->rttpct1 = sp->rttarray[sp->rcvd * percent1 / 100];
                 
        if ((sp->rcvd * percent2) >= 100)
                if (((sp->rcvd * percent2) % 100) == 0)
                        sp->rttpct2 = sp->rttarray[(sp->rcvd * percent2 / 100) - 1];
                else
                        sp->rttpct2 = sp->rttarray[sp->rcvd * percent2 / 100];
        else
                sp->rttpct2 = sp->rttarray[sp->rcvd * percent2 / 100];
        
	if (sp->rcvd == 0) 
		sp->rttmed = 0.0;
}


/* 
** PRINT_STATS -- Prints RTT Statistics 
** ------------------------------------
*/

void
print_stats(sp)
stats *sp;
{
        double rttavg;                  /* average round-trip time */
        double rttstd;                  /* rtt standard deviation */
                
        if (sp->rcvd == 0) {
		sp->rttmin = 0.0;
		sp->rttmax = 0.0;
		sp->rttmed = 0.0;
		rttavg = 0;
		rttstd = 0;
	}

	if (sp->rcvd > 0) 
        {
                rttavg = sp->rttsum / sp->rcvd;
                rttstd = sp->rttssq - (rttavg * sp->rttsum);
		if (sp->rcvd == 1)
	                rttstd = xsqrt(rttstd / sp->rcvd);
		else 
			rttstd = xsqrt(rttstd / (sp->rcvd-1));
        }                

		printf("\n***** Round Trip Statistics of SYN-ACK to %s (Port = %d) ******\n",hisadr->h_name, hisport);
		printf("%d packets transmitted, %d packets received, %2.2f percent packet loss\n",sp->seq_no,sp->rcvd,(((float)sp->seq_no)-(float)sp->rcvd)/((float)sp->seq_no)*100.0);
                printf("round-trip (ms) min/avg/max =");
                printf(" %s", tvprint(sp->rttmin));
                printf("/%s", tvprint((time_t)rttavg));
                printf("/%s", tvprint(sp->rttmax));
                printf(" (std = %s)\n", tvprint((time_t)rttstd));

                statistics(sp,percent1,percent2);
                printf(" (median = %s)\t", tvprint(sp->rttmed));
                printf(" (interquartile range = %s)\n", tvprint(sp->rttiqr));
                if(percent1 != 0)
                        printf(" (%d percentile = %s)\t",percent1,tvprint(sp->rttpct1));
                if(percent2 != 0)
                        printf(" (%d percentile = %s)\n",percent2,tvprint(sp->rttpct2));
                printf("\n");
}


/* 
** XSQRT -- Square root 
** --------------------
*/
 
double
xsqrt(y)
double y;       
{   
        double t, x;

        if (y <= 0)
                return(0);

        x = (y < 1.0) ? 1.0 : y;
        do {
                t = x;
                x = (t + (y/t))/2.0; 
        } while (0 < x && x < t);
                
        return(x);
}


/*
** PRINTSTATS -- PRINTS RTT STATISTICS
** -----------------------------------
*/
void
printstats(sp)
stats *sp;                              /* statistics buffer */
{
 
 	print_stats(sp);
	exit(0);
}


/*
** CALLSTATS -- CALLS RTT STATISTICS
** ---------------------------------
*/
void 
callstats()
{

 	printstats(&synack);

}
  

/* 
** ISNUMERIC -- FINDS OUT IF AN ADDRESS IS NUMERIC
** -----------------------------------------------
*/
int 
isnumeric(address)
char *address;
{

 	if (isdigit(*address))
		return TRUE;
	else 
		return FALSE;

}


/*
** ERROR -- Issue error message to error output
** --------------------------------------------
**
**      Returns:
**              None.
*/               
                
void /*VARARGS1*/
rterror(fmt, a, b, c, d)
char *fmt;                              /* format of message */
char *a, *b, *c, *d;                    /* optional arguments */
{
	(void) fprintf(stderr, "\n");
        (void) fprintf(stderr, fmt, a, b, c, d);
        (void) fprintf(stderr, "\n");
	(void) fprintf(stderr, "\n");
}


/*
** ERR_PRINT -- PRINTS STRING IF ERROR
** -----------------------------------
*/
void 
err_print(string)
char *string;
{
	printf("%s\n\n", string);
	exit(0);
}


/* 
** SOCK_INIT -- OPENS A SOCKET & MAKES A SOCKET NON-BLOCKING
** ---------------------------------------------------------
*/
void sock_init (sock, flag)
int *sock, *flag;
{

	*sock = socket(AF_INET, SOCK_STREAM, 0);
                
        *flag = fcntl(*sock, F_GETFL, 0);
        fcntl(*sock, F_SETFL, *flag | O_NONBLOCK);

} 


/*
** SOCK_CLOSE -- CLOSES A CONNECTION
** ---------------------------------
*/
void sock_close (sock, flag)
int *sock, *flag;
{

	close(*sock);
	fcntl(*sock, F_SETFL, *flag);

}


/*
** CONNECT_START -- INITIATES A CONNECTION
** ---------------------------------------
*/
int connect_start(sock, rset, wset)
int *sock;
fd_set *rset, *wset;
{

 	int stat_connect;

	stat_connect = connect(*sock, (struct sockaddr *)&his_addr, sizeof(struct sockaddr));

        #if(!linux) 
                if (stat_connect == -1) 
                        errno = EINPROGRESS;
	#endif

	if (errno == EINPROGRESS) {
                FD_ZERO(rset);
                FD_ZERO(wset);
                
                FD_SET(*sock, rset);
                FD_SET(*sock, wset);
		
		return 1;
	}

	else return 0;
}	


/*
** SELECT_CONDITION -- ERROR CHECKS FOR SELECT
** -------------------------------------------
*/	
int select_condition(sock, rset, wset)
int *sock;
fd_set *rset, *wset;
{

	int value, nset, nread;
	char *line;
	
	if ( FD_ISSET(*sock, wset) && !FD_ISSET(*sock, rset) )
        	value = 1;
                
	if ( FD_ISSET(*sock, wset) && FD_ISSET(*sock, rset) ) 
		value = 2;
                	
	if( (nset = FD_ISSET(*sock, rset) > 0) ) {
		if( nread = read(*sock, line, 1) < 0) {
                	if (errno == 111)
                        	value = 2;
               		else
                         	value = 3;
             	}	
	}
	
	return value;
}


/*
** SELECT_START -- CALLS THE SELECT FUNCTION
** -----------------------------------------
*/
int select_start(sock, rset, wset, tval)
int *sock;
fd_set *rset,*wset;
struct timeval *tval;
{
	int nselect, nselectval;

	loop: 
		nselect = select(*sock+1, rset, wset, NULL, tval);

	if ( (nselect < 0) && (errno == EINTR) )
		goto loop;
	if ( (nselect < 0) && (errno == EINVAL) ) {
		printf("Error in select value\n");
		exit(0);
	}

	if (nselect == 0) 
		return 0;

	nselectval = select_condition(sock, rset, wset);

	return nselectval;
}


/* 
** CONNECTION_TCP -- THE COMPLETE NON-BLOCKING CONNECTION
** ------------------------------------------------------
*/
void connection_tcp(sock, flag, tvalfresh, rset, wset,ipname, ipnum)
int *sock, *flag;
struct timeval *tvalfresh;
fd_set *rset, *wset;
char *ipname, *ipnum;
{

	struct timeval *tval, *sendtime, tv1, *recvtime, tv2;
	time_t rtt;
	int nconnect, nselect;
	float time;
	int local_seq_no;

	local_seq_no = synack.seq_no;
	

	tval = (struct timeval *) malloc(sizeof(struct timeval *)); 
	
	sock_init(sock, flag);
  
        tval->tv_sec = tvalfresh->tv_sec;
        tval->tv_usec = tvalfresh->tv_usec;
	
        sendtime = &tv1;
        gettimeofday(sendtime, (struct timezone *)NULL);

        nconnect = connect_start(sock, rset, wset);
  
        if (nconnect == 1) {
        
                recvtime = &tv2;
                gettimeofday(recvtime, (struct timezone *)NULL);
                                
                recvtime->tv_sec -= sendtime->tv_sec;
                recvtime->tv_usec -= sendtime->tv_usec;
        
                tval->tv_sec -= recvtime->tv_sec;
                tval->tv_usec = tval->tv_usec - recvtime->tv_usec - 10000;
                                        
                if (tval->tv_usec < 0) {
                        tval->tv_sec -= 1;   
                        tval->tv_usec += 1000000;
                }

      		nselect = select_start(sock, rset, wset, tval);
   
                if (nselect == 0) {
	         	printf("connection for seq no: %d timed out within %s Secs\n",local_seq_no, 
				tvprint(tvalfresh->tv_sec*1000 + tvalfresh->tv_usec/1000));
                        goto jump;
                }
                if (nselect == 1)
                        goto printer;
                if (nselect == 2) {
                        printf("connection refused\n");
                        goto jump;
                }
                if (nselect ==3) {
                        printf("read error\n");
                        goto jump;
                }
                
        } 
                
        else {
                printf("connect not proper\n");
                exit(0);
        }
      
        if (nconnect == 0) {
                /* recvtime = tv2;
                gettimeofday(*recvtime, (struct timezone *)NULL);  */
                printf("nconnect = 0\n");
                exit(0);
        }
        
  
	printer:
        
	recvtime = &tv2;	
	gettimeofday(recvtime, (struct timezone *)NULL);

	rtt = tvsub(recvtime, sendtime);
  	synack.rcvd++;
  	record_stats(&synack,rtt);
        
  	if (ipname != NULL)
        	printf("connected to %s : Seq = %d , RTT = %s ms \n", ipname, local_seq_no, tvprint(rtt));
 	else
        	printf("connected to %s : Seq = %d , RTT = %s ms \n", ipnum, local_seq_no, tvprint(rtt));
                        
  	jump:
                 
  	sock_close (sock, flag);

	free(tval);

}


/*
** CALL_CONNECTION -- CALLS THE PREVIOUS FUNCTION
** ----------------------------------------------
*/ 
void *call_connection () 
{

	connection_tcp(&sockfd[counter], &flags[counter], &tvalfresh, &rset[counter], &wset[counter], ipname, ipnum);

}  


/*
** MAX_THREADS -- CALCULATES THE MAXIMUM NO. ACTIVE THREADS POSSIBLE
** -----------------------------------------------------------------
*/
int max_threads(timeout, interval) 
int timeout;
timer interval;
{
	int threads;
	
	float timeout_usec;
	float total_time;

	timeout_usec = timeout * 1000000;

	total_time = (float) ((interval.sec * 1000000) + interval.usec) ;

	threads = (int) ((timeout_usec/total_time) + 3.0);

	return threads;
}  


/* 
** FINISH -- WAIT FOR OUTSTANDING PACKETS CALLS CALLSTATS()
** --------------------------------------------------------
*/
void finish () 
{

	int i;

	printf("\nWaiting for outstanding packets (if any)..........\n\n");

	sleep(1);

	for (i=0; i <= threads; i++) 
		pthread_join(tid[i], NULL);

	callstats();
}


/* 
** EXIT_INTERRUPT -- KEY SET WHEN INTERRUPTED
** ------------------------------------------
*/ 

void exit_interrupt() 
{

	key = 1;

}


/*
** INTERVAL_BETWEEN_SYNACKS -- TIME BETWEEN SUBSEQUENT SYNs
** --------------------------------------------------------
*/ 
void interval_between_synacks ( timer interval) 
{

	int time;

	time = (interval.usec + (1000000 * interval.sec) );

	usleep(time);
}
