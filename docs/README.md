This repository is built based on Hexagonal Architecture(Ports and Adapter) pattern. The hexagonal architecture was invented by Alistair Cockburn.

> The hexagonal architecture, or ports and adapters architecture, is an architectural pattern used in software design. It aims at creating loosely coupled application components that can be easily connected to their software environment by means of ports and adapters. 


![<img src="hexagonal_architecture_pattern(Ports_land_adapters_pattern)" width="100" height="100"/>](https://user-images.githubusercontent.com/20200714/183576142-3a9b9e3f-b4ce-4ac6-8139-1d2c1dec05db.png)


Let's look at the terminologies now,

**Ports**

> A port can be thought of as a technology-neutral entry point because it chooses the interface through which outside parties can connect with the application, regardless of who or what implements that interface. Similar to how several sorts of devices can communicate with a computer over a USB port if they have an appropriate USB adaptor. Additionally, ports enable the Application to interact with other applications, message brokers, databases, and other external systems and services.

**Adapters**

> An Adapter will initiate interaction with the Application via a Port and a specific technology, such as a REST controller, which allows a client to communicate with the Application. There can be as many Adapters as needed for any single Port without putting the Ports or the Application at risk.

**Application Core**

> The Application is the system's heart; it houses the Application Services that orchestrate the functionality or use cases. The Domain Model, which is the business logic embedded in Aggregates, Entities, and Value Objects, is also included. The Application is represented by a hexagon that receives commands or queries from the Ports and sends requests to other external actors, such as databases, via the Ports. 
When used in conjunction with Domain-Driven Design, the Application, or Hexagon, includes both the Application and Domain layers while leaving the User Interface and Infrastructure layers outside.

**Driving and Driven side**

> The driving (or primary) actors initiate the interaction and are always depicted on the left side. A driving adapter, for example, could be a controller that accepts (user) input and sends it to the Application via a Port.

> The Application "kicks into behaviour" driven (or secondary) actors. For example, the Application may invoke a database Adapter to retrieve a specific data set from persistence.

**Principles of hexagonal architecture**

- Single Responsibility Principle

> The Single Responsibility Principle states that "a component should have only one reason to change." In terms of architecture, this means that if a component has only one reason to change, we don't have to worry about it if the software is changed for any other reason.

- Dependency inversion

> High-level modules should not import anything from low-level modules. Both should depend on abstractions (e.g., interfaces). Abstractions should not depend on details. Details (concrete implementations) should depend on abstractions.


**References:**

- https://alistair.cockburn.us/hexagonal-architecture/

- https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)





